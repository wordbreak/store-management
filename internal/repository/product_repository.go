package repository

import (
	"database/sql"
	"errors"
	"store-management/internal/datasource"
	"store-management/internal/model"
)

type ProductRepository interface {
	CreateProduct(storeId int64, product *model.Product) (int64, error)
	UpdateProduct(product *model.Product) error
	DeleteProduct(storeId, productId int64) error
	FindProductByID(productID int64) (*model.Product, error)
	FindProductsByStoreID(storeID int64, cursor int64, limit int64) ([]*model.Product, error)
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

func (p *productRepositoryImpl) UpdateProduct(product *model.Product) error {
	_, err := p.writer.Exec("UPDATE product SET category = ?, price = ?, cost = ?, name = ?, description = ?, barcode = ?, expiry_date = ?, size = ? WHERE id = ?",
		product.Category, product.Price, product.Cost, product.Name, product.Description, product.Barcode, product.ExpiryDate.Format("2006-01-02 15:04:05"), product.Size, product.ID)
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

func (p *productRepositoryImpl) FindProductByID(productID int64) (*model.Product, error) {
	var product model.Product
	err := p.reader.Get(&product, "SELECT * FROM product WHERE id = ?", productID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, datasource.ErrNoRows
		}
		panic(err)
	}
	return &product, nil
}

func (p *productRepositoryImpl) FindProductsByStoreID(storeID int64, cursor int64, limit int64) ([]*model.Product, error) {
	var products []*model.Product
	err := p.reader.Select(&products, "SELECT * FROM product WHERE store_id = ? AND id < ? ORDER BY id DESC LIMIT ?", storeID, cursor, limit)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, datasource.ErrNoRows
		}
		panic(err)
	}
	return products, nil
}
