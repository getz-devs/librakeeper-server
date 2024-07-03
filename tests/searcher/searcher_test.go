package tests

import (
	searcherv1 "github.com/getz-devs/librakeeper-protos/gen/go/searcher"
	"github.com/getz-devs/librakeeper-server/tests/searcher/suite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math/rand"
	"testing"
	"time"
)

const (
	realIsbn = "978-5-4461-2058-1"
)

func Test1RequestPendingSecondIsFinished(t *testing.T) {
	ctx, st := suite.New(t)

	randomIsbn := GenerateRandomIsbn()
	response, err := st.SearcherServ.SearchByISBN(ctx, &searcherv1.SearchByISBNRequest{Isbn: randomIsbn})
	require.NoError(t, err)
	assert.True(t, response.GetStatus() == searcherv1.SearchByISBNResponse_PROCESSING)

	time.Sleep(5 * time.Second)
	response, err = st.SearcherServ.SearchByISBN(ctx, &searcherv1.SearchByISBNRequest{Isbn: randomIsbn})
	require.NoError(t, err)
	assert.True(t, response.GetStatus() == searcherv1.SearchByISBNResponse_SUCCESS)
}

func TestRealIsbn(t *testing.T) {
	ctx, st := suite.New(t)
	// do-while proccesing
	for {
		response, err := st.SearcherServ.SearchByISBN(ctx, &searcherv1.SearchByISBNRequest{Isbn: realIsbn})
		require.NoError(t, err)
		if response.GetStatus() == searcherv1.SearchByISBNResponse_SUCCESS {
			require.NotEmpty(t, response.GetBooks())
			break
		}
		time.Sleep(2 * time.Second)
	}
}

func GenerateRandomIsbn() string {
	// Generate a random 10-digit ISBN
	isbn := make([]byte, 10)
	for i := 0; i < 9; i++ {
		isbn[i] = byte(rand.Intn(10) + '0')
	}
	isbn[9] = 'X'

	// Convert the byte slice to a string
	return string(isbn)
}
