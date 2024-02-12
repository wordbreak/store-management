package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"store-management/internal/datasource"
	"store-management/internal/model"
	"time"
)

type ProductRepository interface {
	CreateProduct(storeId int64, product *model.Product) (int64, error)
	UpdateProduct(storeId int64, product *model.Product) error
	DeleteProduct(storeId, productId int64) error
	FindProduct(storeId, productId int64) (*model.Product, error)
	FindProductsByStoreID(storeId int64, cursor int64, limit int64) ([]*model.Product, error)
}
type productRepositoryImpl struct {
	writer      datasource.SQL
	reader      datasource.SQL
	transaction datasource.Transaction
}

func NewProductRepository(writer, reader datasource.SQL, transaction datasource.Transaction) ProductRepository {
	return &productRepositoryImpl{
		writer:      writer,
		reader:      reader,
		transaction: transaction,
	}
}

func (p *productRepositoryImpl) CreateProduct(storeId int64, product *model.Product) (int64, error) {
	tx := p.transaction.MustBegin()
	res := tx.MustExec("INSERT INTO product (category, price, cost, name, description, barcode, expiry_date, size) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		product.Category, product.Price, product.Cost, product.Name, product.Description, product.Barcode, product.ExpiryDate.Format("2006-01-02 15:04:05"), product.Size)
	lastInsertId, err := res.LastInsertId()
	if err != nil {
		_ = tx.Rollback()
		return 0, err
	}
	tx.MustExec("INSERT INTO store_product (store_id, product_id) VALUES (?, ?)", storeId, lastInsertId)
	err = tx.Commit()
	if err != nil {
		return 0, err
	}
	return lastInsertId, nil
}

func (p *productRepositoryImpl) UpdateProduct(storeId int64, product *model.Product) error {
	tx := p.transaction.MustBegin()
	var count int
	if err := tx.Get(&count, "SELECT COUNT(id) FROM store_product WHERE store_id = ? AND product_id = ?", storeId, product.ID); err != nil || count == 0 {
		_ = tx.Rollback()
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return err
			}
			panic(err)
		} else {
			return sql.ErrNoRows
		}
	}

	fields := reflect.ValueOf(product).Elem()
	for i := 0; i < fields.NumField(); i++ {
		field := fields.Field(i)
		if field.IsZero() {
			continue
		}
		fieldName := fields.Type().Field(i).Tag.Get("db")
		if fieldName == "ID" {
			continue
		}
		fieldValue := field.Interface()
		if fields.Type().Field(i).Type.String() == "time.Time" {
			if fieldValue.(time.Time).Unix() == 0 {
				continue
			}
			fieldValue = fieldValue.(time.Time).Format("2006-01-02 15:04:05")
		}

		_ = tx.MustExec(fmt.Sprintf("UPDATE product SET %s = ? WHERE id = ?", fieldName), fieldValue, product.ID)
	}

	err := tx.Commit()
	if err != nil {
		panic(err)
	}
	return nil
}

func (p *productRepositoryImpl) DeleteProduct(storeId, productId int64) error {
	tx := p.transaction.MustBegin()
	_, err := tx.Exec("DELETE FROM product WHERE id = ?", productId)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	res := tx.MustExec("DELETE FROM store_product WHERE store_id = ? AND product_id = ?", storeId, productId)
	affectedRows, err := res.RowsAffected()
	if affectedRows == 0 || err != nil {
		_ = tx.Rollback()
		return errors.New("product not found")
	}
	return tx.Commit()
}

func (p *productRepositoryImpl) FindProduct(storeId, productId int64) (*model.Product, error) {
	var count int
	err := p.reader.Get(&count, "SELECT count(id) FROM store_product WHERE store_id = ? AND product_id = ?", storeId, productId)
	if err != nil {
		panic(err)
	}
	if count == 0 {
		return nil, datasource.ErrNoRows
	}

	var product model.Product
	err = p.reader.Get(&product, "SELECT * FROM product WHERE id = ?", productId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, datasource.ErrNoRows
		}
		panic(err)
	}
	return &product, nil
}

func (p *productRepositoryImpl) FindProductsByStoreID(storeId int64, cursor int64, limit int64) ([]*model.Product, error) {
	var products []*model.Product
	err := p.reader.Select(&products, "SELECT * FROM product WHERE store_id = ? AND id < ? ORDER BY id DESC LIMIT ?", storeId, cursor, limit)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, datasource.ErrNoRows
		}
		panic(err)
	}
	return products, nil
}
