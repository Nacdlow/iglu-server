package tokens

import (
	"math/rand"
	"strconv"
	"sync"
	"time"
)

var (
	rnd       *rand.Rand
	validKeys []string
	// We use mutex to prevent race conditions, as this will be accessed from
	// multiple goroutines.
	mut sync.Mutex
)

func init() {
	rnd = rand.New(rand.NewSource(time.Now().UnixNano()))
}

// randIntRange generate a random number within a given range.
func randIntRange(min, max int) int {
	return rand.Intn(max-min) + min
}

// GenerateInviteKey generates a new invite key.
func GenerateInviteKey() (code string) {
	mut.Lock()
	defer mut.Unlock()
	code = strconv.Itoa(randIntRange(1000, 9999))
	validKeys = append(validKeys, code)
	return code
}

// CheckAndConsumeKey checks a given key whether it is valid or not, and
// consumes it if valid. Returns whether the key given was valid.
func CheckAndConsumeKey(key string) bool {
	mut.Lock()
	defer mut.Unlock()
	for i, k := range validKeys {
		if k == key {
			validKeys = append(validKeys[:i], validKeys[i+1:]...)
			return true
		}
	}
	return false
}
