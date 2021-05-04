package main

import (
	"testing"

	"github.com/ipfs/go-ipfs/repo/fsrepo/migrations"
	"github.com/ipfs/go-ipfs/repo/fsrepo/migrations/ipfsfetcher"
)

const peersStr = "/ip4/127.0.0.1/tcp/4001/p2p/12D3KooWGC6TvWhfajpgX6wvJHMYvKpDMXPb3ZnCZ6dMoaMtimQ5,/ip4/127.0.0.1/udp/4001/quic/p2p/12D3KooWGC6TvWhfagifX6wvJHMYvKpDMXPb3ZnCZ6dMoaMtimQ5"

func TestGetMigrationFetcher(t *testing.T) {
	var f migrations.Fetcher
	var err error
	f, err = getMigrationFetcher("ftp://bad.gateway.io", "")
	if err == nil {
		t.Fatal("Expected bad URL scheme error")
	}

	f, err = getMigrationFetcher("ipfs", "")
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := f.(*ipfsfetcher.IpfsFetcher); !ok {
		t.Fatal("expected IpfsFetcher")
	}

	f, err = getMigrationFetcher("http", "")
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := f.(*migrations.HttpFetcher); !ok {
		t.Fatal("expected HttpFetcher")
	}

	f, err = getMigrationFetcher("", "")
	if err != nil {
		t.Fatal(err)
	}
	mf, ok := f.(*migrations.MultiFetcher)
	if !ok {
		t.Fatal("expected MultiFetcher")
	}
	if mf.Len() != 2 {
		t.Fatal("expected 2 fetchers in MultiFetcher")
	}

	f, err = getMigrationFetcher("ipfs,http,some.domain.io", peersStr)
	if err != nil {
		t.Fatal(err)
	}
	mf, ok = f.(*migrations.MultiFetcher)
	if !ok {
		t.Fatal("expected MultiFetcher")
	}
	if mf.Len() != 3 {
		t.Fatal("expected3 fetchers in MultiFetcher")
	}
}

func TestParsePeers(t *testing.T) {
	peers, err := parsePeers(peersStr)
	if err != nil {
		t.Fatal(err)
	}

	if len(peers) != 2 {
		t.Fatal("expected 2 peers, got:", len(peers))
	}

	if peers[0].ID.String() != "12D3KooWGC6TvWhfajpgX6wvJHMYvKpDMXPb3ZnCZ6dMoaMtimQ5" {
		t.Fatal("wrong peer id:", peers[0].ID)
	}
	if peers[0].Addrs[0].String() != "/ip4/127.0.0.1/tcp/4001" {
		t.Fatal("wrong peer addr")
	}
	if peers[1].ID.String() != "12D3KooWGC6TvWhfagifX6wvJHMYvKpDMXPb3ZnCZ6dMoaMtimQ5" {
		t.Fatal("wrong peer id:", peers[1].ID)
	}
	if peers[1].Addrs[0].String() != "/ip4/127.0.0.1/udp/4001/quic" {
		t.Fatal("wrong peer addr")
	}
}
