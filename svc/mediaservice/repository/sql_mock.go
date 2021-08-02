package repository

import (
	"database/sql"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

// NewSQLMock returns mocked database connection
func NewSQLMock() (*sql.DB, sqlmock.Sqlmock, error) {
	return sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
}

// CreateItemMock ...
func CreateItemMock(mock sqlmock.Sqlmock, arg CreateItemParams, expectedErr error) {
	if expectedErr != nil {
		mock.ExpectExec(createItem).
			WillReturnError(expectedErr)
	} else {
		t, _ := time.Parse("Mon Jan 2 15:04:05 -0700 MST 2006", "Mon Jan 2 15:04:05 -0700 MST 2006")
		mock.ExpectQuery(createItem).
			WillReturnRows(
				sqlmock.NewRows([]string{"id", "filename", "filepath", "relation_type", "relation_id", "created_at"}).
					AddRow(arg.ID, arg.Filename, arg.Filepath, arg.RelationType, arg.RelationID, t),
			)
	}
}

// GetItemByIDMock ...
func GetItemByIDMock(mock sqlmock.Sqlmock, arg Item, expectedErr error) {
	if expectedErr != nil {
		mock.ExpectQuery(getItemByID).
			WithArgs(arg.ID).
			WillReturnError(expectedErr)
	} else {
		mock.ExpectQuery(getItemByID).
			WithArgs(arg.ID).
			WillReturnRows(
				sqlmock.NewRows([]string{"id", "filename", "filepath", "relation_type", "relation_id", "created_at"}).
					AddRow(arg.ID, arg.Filename, arg.Filepath, arg.RelationType, arg.RelationID, arg.CreatedAt),
			)
	}
}

// GetItemsListMock ...
func GetItemsListMock(mock sqlmock.Sqlmock, arg GetItemsListParams, items []Item, expectedErr error) {
	if expectedErr != nil {
		mock.ExpectQuery(getItemsList).
			WithArgs(arg.Limit, arg.Offset).
			WillReturnError(expectedErr)
	} else {
		rows := sqlmock.NewRows([]string{"id", "filename", "filepath", "relation_type", "relation_id", "created_at"})
		for _, item := range items {
			rows.AddRow(item.ID, item.Filename, item.Filepath, item.RelationType, item.RelationID, item.CreatedAt)
		}
		mock.ExpectQuery(getItemsList).
			WithArgs(arg.Limit, arg.Offset).
			WillReturnRows(rows)
	}
}
