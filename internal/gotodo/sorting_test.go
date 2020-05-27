package gotodo

import (
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSortbyCreatedDate(t *testing.T) {
	now := time.Now()

	todos := TodoList{
		&Todo{Description: "Index 0", CreationDate: ValidTime(now)},
		&Todo{Description: "Index 1", CreationDate: ValidTime(now.Add(100))},
		&Todo{Description: "Index 2", CreationDate: ValidTime(now.Add(-100))},
	}

	sort.Sort(ByCreatedDate(todos))
	assert.Equal(t, len(todos), 3)
	assert.Equal(t, "Index 2", todos[0].Description)
	assert.Equal(t, "Index 0", todos[1].Description)
	assert.Equal(t, "Index 1", todos[2].Description)
}

func TestSortbyDueDate(t *testing.T) {
	now := time.Now()

	todos := TodoList{
		&Todo{Description: "Index 0", DueDate: ValidTime(now)},
		&Todo{Description: "Index 1", DueDate: ValidTime(now.Add(100))},
		&Todo{Description: "Index 2", DueDate: ValidTime(now.Add(-100))},
		&Todo{Description: "Index 3"},
	}

	sort.Sort(ByDueDate(todos))
	assert.Equal(t, len(todos), 4)
	assert.Equal(t, "Index 2", todos[0].Description)
	assert.Equal(t, "Index 0", todos[1].Description)
	assert.Equal(t, "Index 1", todos[2].Description)
	assert.Equal(t, "Index 3", todos[3].Description)
}

func TestSortByPriority(t *testing.T) {
	todos := TodoList{
		&Todo{Description: "Index 0", Priority: 2},
		&Todo{Description: "Index 1", Priority: 1},
		&Todo{Description: "Index 2", Priority: 0},
	}

	sort.Sort(ByPriority(todos))
	assert.Equal(t, len(todos), 3)
	assert.Equal(t, "Index 1", todos[0].Description)
	assert.Equal(t, "Index 0", todos[1].Description)
	assert.Equal(t, "Index 2", todos[2].Description)
}
