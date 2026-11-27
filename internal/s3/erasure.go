package s3

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Default Erasure Coding layout: 4 data shards + 2 parity shards (tolerates up to 2 drive failures).
const (
	DefaultDataShards   = 4
	DefaultParityShards = 2
)

// GF(2^8) tables with primitive polynomial x^8 + x^4 + x^3 + x^2 + 1 (0x11d).
var (
	gfExp [512]byte
	gfLog [256]byte
)

func init() {
	x := 1
	for i := 0; i < 255; i++ {
		gfExp[i] = byte(x)
		gfLog[x] = byte(i)
		x <<= 1
		if x&0x100 != 0 {
			x ^= 0x11d
		}
	}
	for i := 255; i < 512; i++ {
		gfExp[i] = gfExp[i-255]
	}
}

func gfMul(a, b byte) byte {
	if a == 0 || b == 0 {
		return 0
	}
	return gfExp[int(gfLog[a])+int(gfLog[b])]
}

func gfInv(a byte) byte {
	if a == 0 {
		return 0
	}
	return gfExp[255-int(gfLog[a])]
}

// ErasureSet implements systematic Reed-Solomon erasure coding over GF(2^8).
type ErasureSet struct {
	DataShards   int
	ParityShards int
	TotalShards  int
	encMatrix    [][]byte // TotalShards x DataShards (top DataShards rows = identity)
}

// NewErasureSet creates a systematic Reed-Solomon encoder/decoder.
func NewErasureSet(dataShards, parityShards int) (*ErasureSet, error) {
	if dataShards <= 0 || parityShards <= 0 || dataShards+parityShards > 255 {
		return nil, errors.New("invalid shard configuration")
	}
	total := dataShards + parityShards
	mat := make([][]byte, total)
	for i := 0; i < total; i++ {
		mat[i] = make([]byte, dataShards)
	}
	// Top DataShards rows = Identity
	for i := 0; i < dataShards; i++ {
		mat[i][i] = 1
	}
	// Bottom ParityShards rows = Cauchy matrix 1 / (x_r ^ y_c) which guarantees every square submatrix is invertible
	for r := 0; r < parityShards; r++ {
		xr := byte(dataShards + r)
		for c := 0; c < dataShards; c++ {
			yc := byte(c)
			mat[dataShards+r][c] = gfInv(xr ^ yc)
		}
	}
	return &ErasureSet{
		DataShards:   dataShards,
		ParityShards: parityShards,
		TotalShards:  total,
		encMatrix:    mat,
	}, nil
}

// Encode splits data into DataShards + ParityShards equal-length shards.
func (e *ErasureSet) Encode(data []byte) ([][]byte, error) {
	shardSize := (len(data) + e.DataShards - 1) / e.DataShards
	if shardSize == 0 {
		shardSize = 1
	}
	shards := make([][]byte, e.TotalShards)
	for i := 0; i < e.TotalShards; i++ {
		shards[i] = make([]byte, shardSize)
	}

	// Copy data into data shards (zero-padded at the end)
	for i := 0; i < e.DataShards; i++ {
		start := i * shardSize
		if start < len(data) {
			end := start + shardSize
			if end > len(data) {
				end = len(data)
			}
			copy(shards[i], data[start:end])
		}
	}

	// Compute parity shards: parity[r][b] = sum_c encMatrix[dataShards+r][c] * shards[c][b]
	for r := 0; r < e.ParityShards; r++ {
		row := e.encMatrix[e.DataShards+r]
		dst := shards[e.DataShards+r]
		for c := 0; c < e.DataShards; c++ {
			coeff := row[c]
			if coeff == 0 {
				continue
			}
			src := shards[c]
			for b := 0; b < shardSize; b++ {
				dst[b] ^= gfMul(coeff, src[b])
			}
		}
	}
	return shards, nil
}

// Reconstruct rebuilds missing shards (nil entries in shards) in-place and returns the original data.
func (e *ErasureSet) Reconstruct(shards [][]byte, origSize int64) ([]byte, error) {
	if len(shards) != e.TotalShards {
		return nil, fmt.Errorf("expected %d shards, got %d", e.TotalShards, len(shards))
	}

	shardSize := 0
	validIdx := make([]int, 0, e.DataShards)
	for i, s := range shards {
		if s != nil {
			if shardSize == 0 {
				shardSize = len(s)
			} else if len(s) != shardSize {
				return nil, errors.New("shard size mismatch")
			}
			validIdx = append(validIdx, i)
		}
	}
	if len(validIdx) < e.DataShards {
		return nil, fmt.Errorf("cannot reconstruct: need %d shards, only %d available", e.DataShards, len(validIdx))
	}
	if shardSize == 0 {
		return []byte{}, nil
	}

	// Fast path: all data shards present
	allDataPresent := true
	for i := 0; i < e.DataShards; i++ {
		if shards[i] == nil {
			allDataPresent = false
			break
		}
	}

	if !allDataPresent {
		// Build DataShards x DataShards submatrix from the first DataShards valid rows
		useIdx := validIdx[:e.DataShards]
		sub := make([][]byte, e.DataShards)
		for r, idx := range useIdx {
			sub[r] = make([]byte, e.DataShards)
			copy(sub[r], e.encMatrix[idx])
		}
		inv, err := invertGFMatrix(sub)
		if err != nil {
			return nil, err
		}

		// Reconstruct missing data shards
		for d := 0; d < e.DataShards; d++ {
			if shards[d] != nil {
				continue
			}
			rebuilt := make([]byte, shardSize)
			for c, idx := range useIdx {
				coeff := inv[d][c]
				if coeff == 0 {
					continue
				}
				src := shards[idx]
				for b := 0; b < shardSize; b++ {
					rebuilt[b] ^= gfMul(coeff, src[b])
				}
			}
			shards[d] = rebuilt
		}

		// Recompute any missing parity shards too (for auto-healing)
		for r := 0; r < e.ParityShards; r++ {
			pIdx := e.DataShards + r
			if shards[pIdx] != nil {
				continue
			}
			rebuilt := make([]byte, shardSize)
			row := e.encMatrix[pIdx]
			for c := 0; c < e.DataShards; c++ {
				coeff := row[c]
				if coeff == 0 {
					continue
				}
				src := shards[c]
				for b := 0; b < shardSize; b++ {
					rebuilt[b] ^= gfMul(coeff, src[b])
				}
			}
			shards[pIdx] = rebuilt
		}
	}

	// Join data shards and trim to origSize
	out := make([]byte, 0, int64(e.DataShards)*int64(shardSize))
	for i := 0; i < e.DataShards; i++ {
		out = append(out, shards[i]...)
	}
	if origSize >= 0 && int64(len(out)) > origSize {
		out = out[:origSize]
	}
	return out, nil
}

func invertGFMatrix(m [][]byte) ([][]byte, error) {
	n := len(m)
	aug := make([][]byte, n)
	for i := 0; i < n; i++ {
		aug[i] = make([]byte, 2*n)
		copy(aug[i][:n], m[i])
		aug[i][n+i] = 1
	}

	for col := 0; col < n; col++ {
		pivot := -1
		for r := col; r < n; r++ {
			if aug[r][col] != 0 {
				pivot = r
				break
			}
		}
		if pivot == -1 {
			return nil, errors.New("singular erasure matrix")
		}
		if pivot != col {
			aug[col], aug[pivot] = aug[pivot], aug[col]
		}
		invPivot := gfInv(aug[col][col])
		for c := 0; c < 2*n; c++ {
			aug[col][c] = gfMul(aug[col][c], invPivot)
		}
		for r := 0; r < n; r++ {
			if r == col || aug[r][col] == 0 {
				continue
			}
			factor := aug[r][col]
			for c := 0; c < 2*n; c++ {
				aug[r][c] ^= gfMul(factor, aug[col][c])
			}
		}
	}

	inv := make([][]byte, n)
	for i := 0; i < n; i++ {
		inv[i] = make([]byte, n)
		copy(inv[i], aug[i][n:])
	}
	return inv, nil
}

// DriveDirs returns the configured storage drive directories for erasure-coded shards.
// Uses ORBIT_S3_DRIVES (comma-separated) if set, otherwise 6 virtual drives under StorageDir/drives/.
func (s *S3Engine) DriveDirs() []string {
	if env := strings.TrimSpace(os.Getenv("ORBIT_S3_DRIVES")); env != "" {
		parts := strings.Split(env, ",")
		var dirs []string
		for _, p := range parts {
			p = strings.TrimSpace(p)
			if p != "" {
				dirs = append(dirs, p)
			}
		}
		if len(dirs) >= DefaultDataShards+DefaultParityShards {
			return dirs[:DefaultDataShards+DefaultParityShards]
		}
	}
	dirs := make([]string, DefaultDataShards+DefaultParityShards)
	for i := range dirs {
		dirs[i] = filepath.Join(s.cfg.StorageDir, "drives", fmt.Sprintf("drive-%d", i))
	}
	return dirs
}

// shardFileFormat: [8-byte uint64 origSize][32-byte sha256(shardData)][shardData...]
func packShard(origSize int64, shard []byte) []byte {
	sum := sha256.Sum256(shard)
	buf := make([]byte, 8+32+len(shard))
	binary.BigEndian.PutUint64(buf[:8], uint64(origSize))
	copy(buf[8:40], sum[:])
	copy(buf[40:], shard)
	return buf
}

func unpackShard(raw []byte) (int64, []byte, bool) {
	if len(raw) < 40 {
		return 0, nil, false
	}
	origSize := int64(binary.BigEndian.Uint64(raw[:8]))
	expected := raw[8:40]
	payload := raw[40:]
	actual := sha256.Sum256(payload)
	if !strings.EqualFold(hex.EncodeToString(expected), hex.EncodeToString(actual[:])) {
		return 0, nil, false
	}
	return origSize, payload, true
}

func (s *S3Engine) shardPath(driveDir, bucketName, objectID string, idx int) string {
	return filepath.Join(driveDir, bucketName, fmt.Sprintf("%s.shard.%d", objectID, idx))
}

func (s *S3Engine) writeErasureObject(bucketName, objectID string, data []byte) error {
	es, err := NewErasureSet(DefaultDataShards, DefaultParityShards)
	if err != nil {
		return err
	}
	shards, err := es.Encode(data)
	if err != nil {
		return err
	}
	drives := s.DriveDirs()
	origSize := int64(len(data))
	for i, d := range drives {
		p := s.shardPath(d, bucketName, objectID, i)
		if err := os.MkdirAll(filepath.Dir(p), 0700); err != nil {
			return err
		}
		if err := os.WriteFile(p, packShard(origSize, shards[i]), 0600); err != nil {
			return err
		}
	}
	return nil
}

// readErasureObject reads shards from all drives, reconstructs if up to 2 are missing/corrupt,
// and auto-heals any damaged shards on disk.
func (s *S3Engine) readErasureObject(bucketName, objectID string) (data []byte, healedCount int, err error) {
	es, err := NewErasureSet(DefaultDataShards, DefaultParityShards)
	if err != nil {
		return nil, 0, err
	}
	drives := s.DriveDirs()
	shards := make([][]byte, es.TotalShards)
	missing := make([]bool, es.TotalShards)
	var origSize int64 = -1

	for i, d := range drives {
		raw, rErr := os.ReadFile(s.shardPath(d, bucketName, objectID, i))
		if rErr != nil {
			missing[i] = true
			continue
		}
		sz, payload, ok := unpackShard(raw)
		if !ok {
			missing[i] = true
			continue
		}
		if origSize < 0 {
			origSize = sz
		}
		shards[i] = payload
	}

	if origSize < 0 {
		return nil, 0, errors.New("no valid erasure shards found")
	}

	data, err = es.Reconstruct(shards, origSize)
	if err != nil {
		return nil, 0, err
	}

	// Auto-heal any missing or corrupted shards
	for i, isMissing := range missing {
		if isMissing && shards[i] != nil {
			p := s.shardPath(drives[i], bucketName, objectID, i)
			if mkErr := os.MkdirAll(filepath.Dir(p), 0700); mkErr == nil {
				if wErr := os.WriteFile(p, packShard(origSize, shards[i]), 0600); wErr == nil {
					healedCount++
				}
			}
		}
	}
	return data, healedCount, nil
}

func (s *S3Engine) deleteErasureObject(bucketName, objectID string) {
	for i, d := range s.DriveDirs() {
		_ = os.Remove(s.shardPath(d, bucketName, objectID, i))
	}
}
