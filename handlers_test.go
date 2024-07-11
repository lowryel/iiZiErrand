package rank

import (
	// "net/http/httptest"
	// "strings"
	// "testing"

	// "github.com/gofiber/fiber/v2"
	// "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	// "xorm.io/xorm"
)


// func TestErrandApplicationService(t *testing.T) {
//     // Setup
//     app := fiber.New()
//     // mockDB := setupMockDB() // You'd need to implement this
//     // eas := &Repository{DBConn: &xorm.Engine{}}

//     // Test CreateApplication
//     t.Run("CreateApplication", func(t *testing.T) {
//         app.Post("/applications", eas.CreateApplication)

//         req := httptest.NewRequest("POST", "/applications", strings.NewReader(`{"errand_id":"123","user_id":"456"}`))
//         req.Header.Set("Content-Type", "application/json")
//         resp, _ := app.Test(req)

//         assert.Equal(t, 201, resp.StatusCode)
//         // You can also parse the response body and check the returned data
//     })

//     // Test GetApplicationsByErrandID
//     t.Run("GetApplicationsByErrandID", func(t *testing.T) {
//         app.Get("/applications/errand/:errandID", eas.GetApplicationsByErrandID)

//         req := httptest.NewRequest("GET", "/applications/errand/123", nil)
//         resp, _ := app.Test(req)

//         assert.Equal(t, 200, resp.StatusCode)
//         // Check the response body for expected data
//     })

//     // Test UpdateApplicationStatus
//     t.Run("UpdateApplicationStatus", func(t *testing.T) {
//         app.Patch("/applications/:appID", eas.UpdateApplicationStatus)

//         req := httptest.NewRequest("PATCH", "/applications/1", strings.NewReader(`{"status":"accepted"}`))
//         req.Header.Set("Content-Type", "application/json")
//         resp, _ := app.Test(req)

//         assert.Equal(t, 200, resp.StatusCode)
//         // Check the response body for updated data
//     })
// }




type MockDB struct {
    mock.Mock
}

func (m *MockDB) Insert(beans ...interface{}) (int64, error) {
    args := m.Called(beans)
    return int64(args.Int(0)), args.Error(1)
}

func (m *MockDB) Where(query interface{}, args ...interface{}) *MockDB {
    m.Called(query, args)
    return m
}

func (m *MockDB) Find(beans interface{}) error {
    args := m.Called(beans)
    return args.Error(0)
}

// Implement other necessary methods...

func setupMockDB() *MockDB {
    mockDB := new(MockDB)
    // Setup expectations
    mockDB.On("Insert", mock.Anything).Return(1, nil)
    mockDB.On("Where", mock.Anything, mock.Anything).Return(mockDB)
    mockDB.On("Find", mock.Anything).Return(nil)
    return mockDB
}

