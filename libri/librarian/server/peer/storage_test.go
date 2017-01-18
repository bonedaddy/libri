package peer

import (
	"math/rand"
	"net"
	"testing"
	"time"

	"github.com/drausin/libri/libri/librarian/server/storage"
	"github.com/stretchr/testify/assert"
)

func TestFromStored(t *testing.T) {
	sp := NewTestStoredPeer(rand.New(rand.NewSource(0)), 0)
	p := FromStored(sp)
	AssertPeersEqual(t, sp, p)
}

func TestToStored(t *testing.T) {
	p := NewTestPeer(rand.New(rand.NewSource(0)), 0)
	p.Responses().Success()
	p.Responses().Success()
	p.Responses().Error()
	sp := p.ToStored()
	AssertPeersEqual(t, sp, p)
}

func TestFromStoredAddress(t *testing.T) {
	ip, port := "192.168.1.1", uint32(1000)
	sa := &storage.Address{Ip: ip, Port: port}
	a := fromStoredAddress(sa)
	assert.Equal(t, ip, a.IP.String())
	assert.Equal(t, int(port), a.Port)
}

func TestToStoredAddress(t *testing.T) {
	ip, port := "192.168.1.1", 1000
	a := &net.TCPAddr{IP: net.ParseIP(ip), Port: port}
	sa := toStoredAddress(a)
	assert.Equal(t, ip, sa.Ip)
	assert.Equal(t, uint32(port), sa.Port)
}

func TestFromStoredResponseStats(t *testing.T) {
	now, nQueries, nErrors := time.Now().Unix(), uint64(2), uint64(1)
	from := &storage.Responses{Earliest: now, Latest: now, NQueries: nQueries, NErrors: nErrors}
	to := fromStoredResponseStats(from)
	assert.Equal(t, time.Unix(now, 0).UTC(), to.earliest)
	assert.Equal(t, time.Unix(now, 0).UTC(), to.latest)
	assert.Equal(t, nQueries, to.nQueries)
	assert.Equal(t, nErrors, to.nErrors)
}

func TestToStoredResponseStats(t *testing.T) {
	now, nQueries, nErrors := time.Now().UTC(), uint64(2), uint64(1)
	rs := &responseStats{earliest: now, latest: now, nQueries: nQueries, nErrors: nErrors}
	srs := rs.ToStored()
	assert.Equal(t, now.Unix(), srs.Earliest)
	assert.Equal(t, now.Unix(), srs.Latest)
	assert.Equal(t, nQueries, srs.NQueries)
	assert.Equal(t, nErrors, srs.NErrors)
}