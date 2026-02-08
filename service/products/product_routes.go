package product

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"github.com/ArkaniLoveCoding/Golang-Restfull-Api-MySql/middleware"
	"github.com/ArkaniLoveCoding/Golang-Restfull-Api-MySql/types"
	"github.com/ArkaniLoveCoding/Golang-Restfull-Api-MySql/utils"
)


type HandleRequest struct {
	db types.ProductStore
}

func NewHandlerProduct(db types.ProductStore) *HandleRequest {
	return &HandleRequest{db: db}
}

// testing 

func (h *HandleRequest) CreateNewProductTesting(w http.ResponseWriter, r *http.Request) {

	var request types.PayloadUpdateAndCreate
	if err := utils.DecodeData(r, &request); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Wrong type of the data!", err.Error())
	}


	r.Body = http.MaxBytesReader(w, r.Body, 2 << 20)

	if err := r.ParseMultipartForm(32 << 10); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Failed to settings the multipart form data!", err.Error())
		return
	}


	// define the multipart form value for the request client
	name_product := r.FormValue("name")
	stock_product := r.FormValue("stock")
	category_product := r.FormValue("category")
	expired_product := r.FormValue("expired")
	price_product := r.FormValue("price")

	stock_product_convert, err := strconv.Atoi(stock_product)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Failed to convert the string into a integer!", err.Error())
		return 
	}

	// convert into a json struct
	request.Name = name_product
	request.Price = price_product
	request.Stock = stock_product_convert
	request.Category = category_product
	request.Expired = expired_product

	var validate *validator.Validate

	validate = validator.New()
	if err := validate.Struct(&request); err != nil {
		var errors []string
		for _, errorValidate := range err.(validator.ValidationErrors) {
			errors = append(errors, fmt.Sprintf("Fatal Error ! : %v, %v", errorValidate.Field(), errorValidate.Tag()))
		}
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Failed to get the form file into an image!", err.Error())
		return
	}
	defer file.Close()

	bufff := make([]byte, 512)
	read_file, err := file.Read(bufff)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Failed to read the file from image!", err.Error())
		return
	}
	if read_file == 0 {
		utils.WriteError(w, http.StatusBadRequest, "Failed to read the file from image file!", false)
		return
	}
	content_type := http.DetectContentType(bufff)

	if content_type != "image/png" && content_type != "image/jpeg" {
		utils.WriteError(w, http.StatusBadRequest, "Failed to put your image because the type is not jpg or png!", false)
		return
	}
	file.Seek(0, 0)

	filename := uuid.New().String() + filepath.Ext(header.Filename)
	os.MkdirAll("uploads", os.ModePerm)
	path := filepath.Join("uploads", filename)

	path_file, err := os.Create(path)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Failed to create the os path for these file form!", err.Error())
		return 
	}
	defer path_file.Close()

	io.Copy(path_file, file)

	ctx, cancle := context.WithTimeout(r.Context(), time.Second * 10)
	defer cancle()

	time_created := time.Now().UTC()
	time_updated := time.Now().UTC()

	product := &types.Products{
		Id: uuid.New(),
		Name: request.Name,
		Price: request.Price,
		Stock: request.Stock,
		Image: path,
		Category: request.Category,
		Expired: request.Expired,
		Created_at: time_created,
		Updated_at: time_updated,
	}


	if err := h.db.CreateNewProduct(ctx, product); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Failed to create new data for product clients!", err.Error())
		return
	}

	product_response := types.ProductResponse{
		Id: product.Id,
		Name: product.Name,
		Stock: product.Stock,
		Price: product.Price,
		Category: product.Category,
		Expired: product.Expired,
		Created_at: product.Created_at.Format("2006-01-02"),
		Updated_at: product.Updated_at.Format("2006-01-02"),
	}

	if price_product == "" {
		utils.WriteError(w, http.StatusBadRequest, "Failed to convert from product to product_response!", false)
		return
	}

	utils.WriteSuccess(w, http.StatusOK, "Success to create new Product!", product_response)

}

// controllers

func (h *HandleRequest) CreateProductHandler(w http.ResponseWriter, r *http.Request) {

	role_user, err := middleware.GetValueTokenRole(w, r)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Failed to get role for the user from token!", err.Error())
		return 
	}

	if role_user != "ADMIN" {
		utils.WriteError(w, http.StatusBadRequest, "Cannot access this method because your role is not admin!", false)
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Failed to settings the multipart form data!", err.Error())
		return
	}

	// define the multipart form value for the request client
	name_product := r.FormValue("name")
	stock_product := r.FormValue("stock")
	category_product := r.FormValue("category")
	expired_product := r.FormValue("expired")
	price_product := r.FormValue("price")
	stock_product_convert, err := strconv.Atoi(stock_product)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Failed to convert the string into a integer!", err.Error())
		return 
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Failed to get the form file into an image!", err.Error())
		return
	}
	defer file.Close()

	bufff := make([]byte, 512)
	read_file, err := file.Read(bufff)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Failed to read the file from image!", err.Error())
		return
	}
	if read_file == 0 {
		utils.WriteError(w, http.StatusBadRequest, "Failed to read the file from image file!", false)
		return
	}
	content_type := http.DetectContentType(bufff)

	if content_type != "image/png" && content_type != "image/jpeg" {
		utils.WriteError(w, http.StatusBadRequest, "Failed to put your image because the type is not jpg or png!", false)
		return
	}
	file.Seek(0, 0)

	filename := uuid.New().String() + filepath.Ext(header.Filename)
	os.MkdirAll("uploads", os.ModePerm)
	path := filepath.Join("uploads", filename)

	path_file, err := os.Create(path)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Failed to create the os path for these file form!", err.Error())
		return 
	}
	defer path_file.Close()

	io.Copy(path_file, file)

	ctx, cancle := context.WithTimeout(r.Context(), time.Second * 10)
	defer cancle()

	time_created := time.Now().UTC().Format("2006-01-02")
	time_updated := time.Now().UTC().Format("2006-01-02")

	products := types.PayloadUpdateAndCreate{
		Id: uuid.New(),
		Name: name_product,
		Stock: stock_product_convert,
		Image: path,
		Price: price_product,
		Category: category_product,
		Expired: expired_product,
		Created_at: time_created,
		Updated_at: time_updated,
	}	

	products_types := &types.Products{
		Id: products.Id,
		Name: products.Name,
		Stock: products.Stock,
		Image: products.Image,
		Price: products.Price,
		Category: products.Category,
		Expired: products.Expired,
		Created_at: time.Now().UTC(),
		Updated_at: time.Now().UTC(),
	}

	var validate *validator.Validate
	validate = validator.New()
	if err := validate.Struct(&products); err != nil {
		var errors []string
		for _, errorValidate := range err.(validator.ValidationErrors) {
			errors = append(errors, fmt.Sprintf("Fatal Error ! : %v, %v", errorValidate.Field(), errorValidate.Tag()))
		}
	}

	if err := h.db.CreateNewProduct(ctx, products_types); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Failed to create a new product!", err.Error())
		return
	}

	product_response := types.ProductResponse{
		Id: products_types.Id,
		Name: products_types.Name,
		Stock: products_types.Stock,
		Image: products_types.Image,
		Price: products_types.Price,
		Category: products_types.Category,
		Expired: products_types.Expired,
		Created_at: products.Created_at,
		Updated_at: products.Updated_at,
	}

	utils.WriteSuccess(w, http.StatusOK, "Successfully to create new Product!", product_response)

}

func (h *HandleRequest) GetProductByIDHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]

	uuid_parse_id, err := uuid.Parse(id)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Failed to convert into uuid!", err.Error())
		return
	}

	if id == "" {
		utils.WriteError(w, http.StatusBadRequest, "Cannot find the params!", false)
		return 
	}

	products, err := h.db.GetProductByID(uuid_parse_id)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Cannot find the data from id!", err.Error())
		return 
	}

	if products == nil {
		utils.WriteError(w, http.StatusBadRequest, "The id that you want to find is nil!", false)
		return 
	}

	product_response := &types.Products{
		Id: products.Id,
		Name: products.Name,
		Stock: products.Stock,
		Image: products.Image,
		Price: products.Price,
		Expired: products.Expired,
		Category: products.Category,
		Created_at: products.Created_at,
		Updated_at: products.Updated_at,
	}

	utils.WriteSuccess(w, http.StatusOK, "Successfully!", product_response)

}

func (h *HandleRequest) GetAllProduct(w http.ResponseWriter, r *http.Request) {

	products, err := h.db.GetAllProduct()
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Failed to find all data of product!", err.Error())
		return 
	}

	if products == nil {
		utils.WriteError(w, http.StatusBadRequest, "Failed to get data of products because the data is nil!", false)
		return
	}

	utils.WriteSuccess(w, http.StatusAccepted, "Successfully to get all data of product", products)

}

func (h *HandleRequest) DeleteProduct(w http.ResponseWriter, r *http.Request) {

	middleware_checking_role, err := middleware.GetValueTokenRole(w, r)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Failed to get role from jwt token!", err.Error())
		return 
	}

	if middleware_checking_role != "ADMIN" {
		utils.WriteError(w, http.StatusBadRequest, "Only admin can delete one of the data of the product db!", false)
		return
	}

	vars_id := mux.Vars(r)
	id := vars_id["id"]

	if id == "" {
		utils.WriteError(w, http.StatusBadRequest, "the id that you want to be a params is nil!", false)
		return 
	}

	uuid_parse_id, err := uuid.Parse(id)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Failed to convert the data to uuid!", err.Error())
		return
	}

	users, err := h.db.GetProductByID(uuid_parse_id)

	if users == nil {
		utils.WriteError(w, http.StatusBadRequest, "Failed to get users from db use id at param !", false)
		return 
	}

	ctx, cancle := context.WithTimeout(context.Background(), time.Second * 10)
	defer cancle()

	if err := h.db.DeleteProductsOnlyAdmin(uuid_parse_id, ctx); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Failed to delete the product!", err.Error())
		return 
	}

	utils.WriteSuccess(w, http.StatusOK, "Successfully to delete the one of data of the product db!", true)
	
}

func (h *HandleRequest) UpdateProductsOnlyAdmin(w http.ResponseWriter, r *http.Request) {

	

}

func (h *HandleRequest) SearchManyProductsRoutes(w http.ResponseWriter, r *http.Request) {

	keyword := r.URL.Query()
	keyword_user := keyword.Get("products")
	keyword_page := keyword.Get("page")
	keyword_limit := keyword.Get("limit")

	page, err := strconv.Atoi(keyword_page)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Failed to decode the params of page!", err.Error())
		return
	}

	limit, err := strconv.Atoi(keyword_limit)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Failed to get the params of the limit!", err.Error())
	}

	if page < 1 {
		page = 1
	}

	offset := (page - 1) * limit

	ctx, cancle := context.WithTimeout(r.Context(), time.Second * 10)
	defer cancle()
	
	products, err := h.db.SearchManyProducts(
		ctx,
		keyword_user,
		offset,
		limit,
	)

	if products == nil {
		utils.WriteError(w, http.StatusBadRequest, "Failed to find the data as you want!", false)
		return
	}

	utils.WriteSuccess(w, http.StatusAccepted, "Successfully to get data from db!", products)

}
