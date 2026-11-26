package tlsengine

import (
	"bufio"
	"context"
	"errors"
	"net"
	"net/http"
	"net/netip"
	"strings"
	"sync"
	"time"
)

const cloudflareIPListV4 = "https://www.cloudflare.com/ips-v4"
const cloudflareIPListV6 = "https://www.cloudflare.com/ips-v6"

// Detection says whether a domain resolves to Cloudflare, and which provider to recommend.
type Detection struct {
	Domain      string   `json:"domain"`
	IPs         []string `json:"ips"`
	Cloudflare  bool     `json:"cloudflare"`
	Recommended string   `json:"recommended"`
}

var (
	cfMu       sync.Mutex
	cfPrefixes []netip.Prefix
	cfFetched  time.Time
)

// Detect resolves the domain and checks each address against Cloudflare's published ranges.
// If the Cloudflare list cannot be fetched, it reports an error instead of guessing.
func Detect(ctx context.Context, domain string) (*Detection, error) {
	d, err := NormalizeDomain(domain)
	if err != nil {
		return nil, err
	}
	addrs, err := net.DefaultResolver.LookupHost(ctx, d)
	if err != nil {
		return nil, errors.New("domain does not resolve yet: point its DNS record at this server first")
	}
	prefixes, err := cloudflarePrefixes(ctx)
	if err != nil {
		return nil, err
	}
	out := &Detection{Domain: d, IPs: addrs, Recommended: ProviderLetsEncrypt}
	for _, a := range addrs {
		ip, err := netip.ParseAddr(a)
		if err != nil {
			continue
		}
		for _, p := range prefixes {
			if p.Contains(ip) {
				out.Cloudflare = true
			}
		}
	}
	if out.Cloudflare {
		out.Recommended = ProviderCloudflareOrigin
	}
	return out, nil
}

// cloudflarePrefixes returns Cloudflare's IPv4 and IPv6 ranges, cached for a day.
func cloudflarePrefixes(ctx context.Context) ([]netip.Prefix, error) {
	cfMu.Lock()
	defer cfMu.Unlock()
	if len(cfPrefixes) > 0 && time.Since(cfFetched) < 24*time.Hour {
		return cfPrefixes, nil
	}
	var all []netip.Prefix
	for _, u := range []string{cloudflareIPListV4, cloudflareIPListV6} {
		ps, err := fetchPrefixes(ctx, u)
		if err != nil {
			return nil, errors.New("could not load Cloudflare IP ranges: " + err.Error())
		}
		all = append(all, ps...)
	}
	cfPrefixes, cfFetched = all, time.Now()
	return all, nil
}

func fetchPrefixes(ctx context.Context, url string) ([]netip.Prefix, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, errors.New(url + " returned " + res.Status)
	}
	var out []netip.Prefix
	sc := bufio.NewScanner(res.Body)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		p, err := netip.ParsePrefix(line)
		if err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, sc.Err()
}
