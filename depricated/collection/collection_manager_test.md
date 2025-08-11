package collection

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/mahdi-cpp/api-go-pkg/asset"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestItem implements CollectionItem for testing
type TestItem struct {
	ID               int
	CreationDate     time.Time
	ModificationDate time.Time
	Name             string
}

func (i *TestItem) GetID() int                      { return i.ID }
func (i *TestItem) SetID(id int)                    { i.ID = id }
func (i *TestItem) SetCreationDate(t time.Time)     { i.CreationDate = t }
func (i *TestItem) SetModificationDate(t time.Time) { i.ModificationDate = t }
func (i *TestItem) GetCreationDate() time.Time      { return i.CreationDate }
func (i *TestItem) GetModificationDate() time.Time  { return i.ModificationDate }

func setupTestManager(t *testing.T) (*Manager[*TestItem], string) {
	t.Helper()

	// Create a temporary file for testing
	tmpFile, err := os.CreateTemp("", "collection_test")
	require.NoError(t, err)
	tmpPath := tmpFile.Name()
	tmpFile.Close()

	manager, err := NewCollectionManager[*TestItem](tmpPath)
	require.NoError(t, err)

	return manager, tmpPath
}

func cleanupTestManager(t *testing.T, path string) {
	t.Helper()
	os.Remove(path)
}

func TestNewCollectionManager(t *testing.T) {
	manager, path := setupTestManager(t)
	defer cleanupTestManager(t, path)

	assert.NotNil(t, manager)
	assert.NotNil(t, manager.items)
	assert.NotNil(t, manager.metadata)
	assert.NotNil(t, manager.ItemAssets)
}

func TestCreateItem(t *testing.T) {
	manager, path := setupTestManager(t)
	defer cleanupTestManager(t, path)

	item := &TestItem{Name: "Test Item"}
	createdItem, err := manager.Create(item)
	require.NoError(t, err)

	assert.Equal(t, 1, createdItem.GetID())
	assert.False(t, createdItem.GetCreationDate().IsZero())
	assert.False(t, createdItem.GetModificationDate().IsZero())
	assert.Equal(t, "Test Item", createdItem.Name)

	// Verify the item was added to the registry
	retrievedItem, err := manager.Get(1)
	require.NoError(t, err)
	assert.Equal(t, createdItem.GetID(), retrievedItem.GetID())
}

func TestUpdateItem(t *testing.T) {
	manager, path := setupTestManager(t)
	defer cleanupTestManager(t, path)

	// Create initial item
	item := &TestItem{Name: "Original"}
	createdItem, err := manager.Create(item)
	require.NoError(t, err)

	// Update the item
	createdItem.Name = "Updated"
	updatedItem, err := manager.Update(createdItem)
	require.NoError(t, err)

	assert.Equal(t, "Updated", updatedItem.Name)
	assert.Equal(t, createdItem.GetID(), updatedItem.GetID())
	assert.False(t, updatedItem.GetModificationDate().Equal(createdItem.GetModificationDate()))

	// Verify the update persisted
	retrievedItem, err := manager.Get(createdItem.GetID())
	require.NoError(t, err)
	assert.Equal(t, "Updated", retrievedItem.Name)
}

func TestDeleteItem(t *testing.T) {
	manager, path := setupTestManager(t)
	defer cleanupTestManager(t, path)

	// Create an item to delete
	item := &TestItem{Name: "To Delete"}
	createdItem, err := manager.Create(item)
	require.NoError(t, err)

	// Delete the item
	err = manager.Delete(createdItem.GetID())
	require.NoError(t, err)

	// Verify it's gone
	_, err = manager.Get(createdItem.GetID())
	assert.True(t, errors.Is(err, errors.New("item not found")))
}

func TestGetAllItems(t *testing.T) {
	manager, path := setupTestManager(t)
	defer cleanupTestManager(t, path)

	// Create several items
	item1 := &TestItem{Name: "Item 1"}
	item2 := &TestItem{Name: "Item 2"}
	item3 := &TestItem{Name: "Item 3"}

	_, err := manager.Create(item1)
	require.NoError(t, err)
	_, err = manager.Create(item2)
	require.NoError(t, err)
	_, err = manager.Create(item3)
	require.NoError(t, err)

	// Get all items
	items, err := manager.GetAll()
	require.NoError(t, err)
	assert.Len(t, items, 3)
}

func TestGetListWithFilter(t *testing.T) {
	manager, path := setupTestManager(t)
	defer cleanupTestManager(t, path)

	// Create test items
	item1 := &TestItem{Name: "Apple"}
	item2 := &TestItem{Name: "Banana"}
	item3 := &TestItem{Name: "Apple Pie"}

	_, err := manager.Create(item1)
	require.NoError(t, err)
	_, err = manager.Create(item2)
	require.NoError(t, err)
	_, err = manager.Create(item3)
	require.NoError(t, err)

	// Filter for items with "Apple" in the name
	filtered, err := manager.GetList(func(item *TestItem) bool {
		return strings.Contains(item.Name, "Apple")
	})
	require.NoError(t, err)
	assert.Len(t, filtered, 2)
}

func TestSortItems(t *testing.T) {
	manager, path := setupTestManager(t)
	defer cleanupTestManager(t, path)

	now := time.Now()

	// Create test items with specific dates for reliable sorting
	item1 := &TestItem{ID: 3, CreationDate: now.Add(-2 * time.Hour), ModificationDate: now.Add(-1 * time.Hour)}
	item2 := &TestItem{ID: 1, CreationDate: now.Add(-1 * time.Hour), ModificationDate: now.Add(-3 * time.Hour)}
	item3 := &TestItem{ID: 2, CreationDate: now, ModificationDate: now.Add(-2 * time.Hour)}

	// Add directly to registry for testing (bypassing Create to control dates)
	manager.items.Register(strconv.Itoa(item1.GetID()), item1)
	manager.items.Register(strconv.Itoa(item2.GetID()), item2)
	manager.items.Register(strconv.Itoa(item3.GetID()), item3)

	// Get all items
	items, err := manager.GetAll()
	require.NoError(t, err)
	require.Len(t, items, 3)

	// Test sorting by ID ascending
	sorted := manager.SortItems(items, SortOptions{SortBy: "id", SortOrder: "asc"})
	assert.Equal(t, 1, sorted[0].GetID())
	assert.Equal(t, 2, sorted[1].GetID())
	assert.Equal(t, 3, sorted[2].GetID())

	// Test sorting by ID descending
	sorted = manager.SortItems(items, SortOptions{SortBy: "id", SortOrder: "desc"})
	assert.Equal(t, 3, sorted[0].GetID())
	assert.Equal(t, 2, sorted[1].GetID())
	assert.Equal(t, 1, sorted[2].GetID())

	// Test sorting by creation date ascending
	sorted = manager.SortItems(items, SortOptions{SortBy: "creationDate", SortOrder: "asc"})
	assert.Equal(t, 3, sorted[0].GetID()) // Oldest creation date
	assert.Equal(t, 1, sorted[1].GetID())
	assert.Equal(t, 2, sorted[2].GetID()) // Newest creation date

	// Test sorting by modification date descending
	sorted = manager.SortItems(items, SortOptions{SortBy: "modificationDate", SortOrder: "desc"})
	assert.Equal(t, 1, sorted[0].GetID()) // Most recent modification
	assert.Equal(t, 3, sorted[1].GetID())
	assert.Equal(t, 2, sorted[2].GetID()) // Oldest modification
}

func TestGetItemAssets(t *testing.T) {
	manager, path := setupTestManager(t)
	defer cleanupTestManager(t, path)

	// Create an item
	item := &TestItem{Name: "With Assets"}
	createdItem, err := manager.Create(item)
	require.NoError(t, err)

	// Add assets to the item
	assets := []*asset.PHAsset{
		{ID: 1, Filename: "Asset 1"},
		{ID: 2, Filename: "Asset 2"},
	}
	manager.ItemAssets[createdItem.GetID()] = assets

	// Retrieve the assets
	retrievedAssets, err := manager.GetItemAssets(createdItem.GetID())
	require.NoError(t, err)
	assert.Len(t, retrievedAssets, 2)
	assert.Equal(t, "Asset 1", retrievedAssets[0].Filename)
}

func TestGetSortedList(t *testing.T) {
	manager, path := setupTestManager(t)
	defer cleanupTestManager(t, path)

	// Create test items
	item1 := &TestItem{ID: 3, Name: "C Item"}
	item2 := &TestItem{ID: 1, Name: "A Item"}
	item3 := &TestItem{ID: 2, Name: "B Item"}

	_, err := manager.Create(item1)
	require.NoError(t, err)
	_, err = manager.Create(item2)
	require.NoError(t, err)
	_, err = manager.Create(item3)
	require.NoError(t, err)

	// Get sorted list by ID ascending
	sorted, err := manager.GetSortedList(nil, "id", "asc")
	require.NoError(t, err)
	assert.Equal(t, 1, sorted[0].GetID())
	assert.Equal(t, 2, sorted[1].GetID())
	assert.Equal(t, 3, sorted[2].GetID())

	// Get sorted list with filter
	filteredSorted, err := manager.GetSortedList(
		func(item *TestItem) bool { return item.GetID() > 1 },
		"id",
		"asc",
	)
	require.NoError(t, err)
	assert.Len(t, filteredSorted, 2)
	assert.Equal(t, 2, filteredSorted[0].GetID())
	assert.Equal(t, 3, filteredSorted[1].GetID())
}

func TestGetAllSorted(t *testing.T) {
	manager, path := setupTestManager(t)
	defer cleanupTestManager(t, path)

	// Create test items
	item1 := &TestItem{ID: 3, Name: "C Item"}
	item2 := &TestItem{ID: 1, Name: "A Item"}
	item3 := &TestItem{ID: 2, Name: "B Item"}

	_, err := manager.Create(item1)
	require.NoError(t, err)
	_, err = manager.Create(item2)
	require.NoError(t, err)
	_, err = manager.Create(item3)
	require.NoError(t, err)

	// Get all sorted by ID descending
	sorted, err := manager.GetAllSorted("id", "desc")
	require.NoError(t, err)
	assert.Equal(t, 3, sorted[0].GetID())
	assert.Equal(t, 2, sorted[1].GetID())
	assert.Equal(t, 1, sorted[2].GetID())
}

func TestConcurrentAccess(t *testing.T) {
	manager, path := setupTestManager(t)
	defer cleanupTestManager(t, path)

	// Test concurrent creates
	var wg sync.WaitGroup
	count := 10

	wg.Add(count)
	for i := 0; i < count; i++ {
		go func() {
			defer wg.Done()
			_, err := manager.Create(&TestItem{})
			assert.NoError(t, err)
		}()
	}
	wg.Wait()

	// Verify all items were created with unique IDs
	items, err := manager.GetAll()
	require.NoError(t, err)
	assert.Len(t, items, count)

	// Check for unique IDs
	idSet := make(map[int]bool)
	for _, item := range items {
		idSet[item.GetID()] = true
	}
	assert.Len(t, idSet, count)
}
