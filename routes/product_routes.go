package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hoffax/prodrest/services"
)

func (h *Handlers) RegisterProductRoutes() {
	g := h.app.Group("/products")
	g.Get("/", h.getAllProducts)
	g.Get("/:id", h.getProductById)
	g.Post("/", h.createProduct)
	g.Put("/:id", h.updateProduct)

	h.app.Get("/check_barcode/:barcode", h.checkBarcode)
}

type GetAllProductsQuery struct {
	StatusOptions []string `query:"status"`
	Search        string   `query:"search"`
	Limit         int      `query:"limit"`
	Offset        int      `query:"offset"`
}

func (h *Handlers) getAllProducts(c *fiber.Ctx) error {
	params := new(GetAllProductsQuery)
	if err := c.QueryParser(params); err != nil {
		return err
	}

	if params.Limit == 0 {
		params.Limit = 10
	}

	products, err := h.sm.FetchProducts(c.Context(), &services.FetchProductsParams{
		StatusOptions: params.StatusOptions,
		Search:        params.Search,
		Limit:         params.Limit,
		Offset:        params.Offset,
	})
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(products)
}

func (h *Handlers) getProductById(c *fiber.Ctx) error {
	productId, err := h.getIdParam(c)
	if err != nil {
		return err
	}

	product, err := h.sm.FetchProductById(c.Context(), productId)
	if err != nil {
		return err
	}

	if product == nil {
		return c.Status(fiber.StatusNotFound).Send([]byte{})
	}

	return c.Status(fiber.StatusOK).JSON(product)
}

type CreateProductBody struct {
	Name             string `json:"name"`
	Barcode          string `json:"barcode"`
	Unit             string `json:"unit" `
	BatchControl     bool   `json:"batchControl"`
	ConversionFactor int    `json:"conversionFactor" `
}

func (h *Handlers) createProduct(c *fiber.Ctx) error {
	params := new(CreateProductBody)
	if err := c.BodyParser(params); err != nil {
		return err
	}

	product, err := h.sm.CreateProduct(c.Context(), &services.CreateProductParams{
		Name:             params.Name,
		Barcode:          params.Barcode,
		Unit:             params.Unit,
		BatchControl:     params.BatchControl,
		ConversionFactor: params.ConversionFactor,
	})
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(product)
}

type UpdateProductBody struct {
	Status           string `json:"status"`
	Name             string `json:"name"`
	Barcode          string `json:"barcode"`
	Unit             string `json:"unit" `
	BatchControl     bool   `json:"batchControl"`
	ConversionFactor int    `json:"conversionFactor" `
}

func (h *Handlers) updateProduct(c *fiber.Ctx) error {
	productId, err := h.getIdParam(c)
	if err != nil {
		return err
	}

	params := new(UpdateProductBody)
	if err := c.BodyParser(params); err != nil {
		return err
	}

	product, err := h.sm.UpdateProduct(c.Context(), &services.UpdateProductParams{
		ID:               productId,
		Status:           params.Status,
		Name:             params.Name,
		Barcode:          params.Barcode,
		Unit:             params.Unit,
		BatchControl:     params.BatchControl,
		ConversionFactor: params.ConversionFactor,
	})
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(product)
}

func (h *Handlers) checkBarcode(c *fiber.Ctx) error {
	barcode := c.Params("barcode")

	product, err := h.sm.FetchProductByBarcode(c.Context(), barcode)
	if err != nil {
		return err
	}

	if product == nil {
		return c.Status(fiber.StatusNotFound).Send([]byte{})
	}

	return c.Status(fiber.StatusOK).JSON(product)
}
