package infra

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"github.com/lil-oren/rest/internal/dependency"
	"github.com/lil-oren/rest/internal/handler/resthandler"
	"github.com/lil-oren/rest/internal/middleware"
	"github.com/lil-oren/rest/internal/repository"
	"github.com/lil-oren/rest/internal/shared"
	"github.com/lil-oren/rest/internal/usecase"
)

type (
	server struct {
		r            *gin.Engine
		v            *validator.Validate
		repositories repositories
		usecases     usecases
		cfg          dependency.Config
	}

	repositories struct {
		exampleRepository        repository.ExampleRepository
		accountRepository        repository.AccountRepository
		accountAddressRepository repository.AccountAddressRepository
		productRepository        repository.ProductRepository
		productVariantRepository repository.ProductVariantRepository
		productMediaRepository   repository.ProductMediaRepository
		shopRepository           repository.ShopRepository
		variantTypeRepository    repository.VariantTypeRepository
		variantGroupRepository   repository.VariantGroupRepository
		provinceRepository       repository.ProvinceRepository
		districtRepository       repository.DistrictRepository
		shopCourierRepository    repository.ShopCourierRepository
		cacheRepository          repository.CacheRepository
		cartRepository           repository.CartRepository
		walletRepository         repository.WalletRepository
		orderRepository          repository.OrderRepository
		courierRepository        repository.CourierRepository
		rajaOngkirRepository     repository.RajaOngkirRepository
		changedEmailRepository   repository.ChangedEmailRepository
		transactionRepository    repository.TransactionRepository
		wishlistRepository       repository.WishlistRepository
		sellerPageRepository     repository.SellerPageRepository
		categoryRepository       repository.CategoryRepository
		reviewRepository         repository.ReviewRepository
		promotionRepository      repository.PromotionRepository
		orderDetailRepository    repository.OrderDetailRepository
	}

	usecases struct {
		exampleUsecase        usecase.ExampleUsecase
		authUsecase           usecase.AuthUsecase
		homepageUsecase       usecase.HomepageUsecase
		accountAddressUsecase usecase.ProfileUsecase
		productPageUsecase    usecase.ProductPageUsecase
		shopUsecase           usecase.ShopUsecase
		dropdownUsecase       usecase.DropdownUsecase
		cartUsecase           usecase.CartUsecase
		walletUsecase         usecase.WalletUsecase
		orderSellerUsecase    usecase.OrderSellerUsecase
		orderUsecase          usecase.OrderUsecase
		checkoutUsecase       usecase.CheckoutUsecase
		discoveryUsecase      usecase.DiscoveryUsecase
		sellerPageUsecase     usecase.SellerPageUsecase
		wishlistUseCase       usecase.WishlistUseCase
		reviewUsecase         usecase.ReviewUsecase
		promotionUsecase      usecase.PromotionUsecase
	}
)

func (s *server) initRepository(db *sqlx.DB, rd *redis.Client, cfg dependency.Config) {

	s.repositories.exampleRepository = repository.NewExampleRepository()
	s.repositories.categoryRepository = repository.NewCategoryRepository(db)
	s.repositories.transactionRepository = repository.NewTransactionRepository(db)
	s.repositories.walletRepository = repository.NewWalletRepository(db, s.repositories.transactionRepository)
	s.repositories.accountRepository = repository.NewAccountRepository(db, s.repositories.walletRepository)
	s.repositories.accountAddressRepository = repository.NewAccountAddressRepository(db)
	s.repositories.productRepository = repository.NewProductRepository(db)
	s.repositories.shopRepository = repository.NewShopRepository(db)
	s.repositories.productVariantRepository = repository.NewProductVariantRepository(db)
	s.repositories.productMediaRepository = repository.NewProductMediaRepository(db)
	s.repositories.shopRepository = repository.NewShopRepository(db)
	s.repositories.variantTypeRepository = repository.NewVariantTypeRepository(db)
	s.repositories.variantGroupRepository = repository.NewVariantGroupRepository(db)
	s.repositories.provinceRepository = repository.NewProvinceRepository(db)
	s.repositories.districtRepository = repository.NewDistrictRepository(db)
	s.repositories.shopCourierRepository = repository.NewShopCourierRepository(db)
	s.repositories.cartRepository = repository.NewCartRepository(db)
	s.repositories.cacheRepository = repository.NewCacheRepository(rd, s.cfg)
	s.repositories.cartRepository = repository.NewCartRepository(db)
	s.repositories.orderRepository = repository.NewOrderRepository(db, s.repositories.transactionRepository)
	s.repositories.courierRepository = repository.NewCourierRepository(db)
	s.repositories.rajaOngkirRepository = repository.NewRajaOngkirRepository(cfg)
	s.repositories.changedEmailRepository = repository.NewChangedEmailRepository(db)
	s.repositories.sellerPageRepository = repository.NewSellerPageRepository(db)
	s.repositories.wishlistRepository = repository.NewWishlistRepository(db)
	s.repositories.reviewRepository = repository.NewReviewRepository(db)
	s.repositories.promotionRepository = repository.NewPromotionRepository(db)
	s.repositories.orderDetailRepository = repository.NewOrderDetailRepository(db)
}

func (s *server) initUsecase(rd *redis.Client) {
	s.usecases.exampleUsecase = usecase.NewExampleRepository()
	s.usecases.authUsecase = usecase.NewAuthUsecase(
		s.repositories.accountRepository,
		s.repositories.cacheRepository,
		s.repositories.cartRepository,
		s.repositories.walletRepository,
		s.repositories.changedEmailRepository,
		s.repositories.shopRepository,
		s.cfg,
	)
	s.usecases.homepageUsecase = usecase.NewHomepageUsecase(
		s.repositories.productRepository,
		s.repositories.cartRepository,
		s.repositories.reviewRepository,
		s.repositories.categoryRepository,
		s.repositories.cacheRepository,
	)
	s.usecases.accountAddressUsecase = usecase.NewProfileUsecase(
		s.repositories.accountAddressRepository,
		s.repositories.accountRepository,
		s.repositories.districtRepository,
		s.repositories.provinceRepository,
	)
	s.usecases.shopUsecase = usecase.NewShopUsecase(
		s.repositories.shopRepository,
		s.repositories.accountAddressRepository,
		s.repositories.shopCourierRepository,
		s.repositories.walletRepository,
		s.repositories.productRepository,
	)
	s.usecases.productPageUsecase = usecase.NewProductPageUsecase(
		s.repositories.productRepository,
		s.repositories.shopRepository,
		s.repositories.productVariantRepository,
		s.repositories.productMediaRepository,
		s.repositories.variantTypeRepository,
		s.repositories.variantGroupRepository,
		s.repositories.wishlistRepository,
		s.repositories.reviewRepository,
		s.repositories.orderDetailRepository,
	)
	s.usecases.cartUsecase = usecase.NewCartUsecase(s.repositories.cartRepository, s.repositories.productVariantRepository)
	s.usecases.dropdownUsecase = usecase.NewDropdownUsecase(
		s.repositories.provinceRepository,
		s.repositories.districtRepository,
		s.repositories.shopCourierRepository,
		s.repositories.categoryRepository,
	)
	s.usecases.walletUsecase = usecase.NewWalletUsecase(
		s.repositories.walletRepository,
		s.cfg,
		s.repositories.transactionRepository,
		s.repositories.accountRepository,
	)
	s.usecases.orderUsecase = usecase.NewOrderUsecase(
		s.repositories.orderRepository,
		s.repositories.cartRepository,
		s.repositories.walletRepository,
		s.repositories.rajaOngkirRepository,
		s.repositories.accountAddressRepository,
		s.repositories.courierRepository,
		s.repositories.productRepository,
		s.repositories.transactionRepository,
		s.repositories.promotionRepository,
	)
	s.usecases.orderSellerUsecase = usecase.NewOrderSellerUsecase(
		s.repositories.orderRepository,
		s.repositories.transactionRepository,
		s.repositories.walletRepository,
	)
	s.usecases.checkoutUsecase = usecase.NewCheckoutUsecase(
		s.repositories.cartRepository,
		s.repositories.accountAddressRepository,
		s.repositories.courierRepository,
		s.repositories.productRepository,
		s.repositories.rajaOngkirRepository,
		s.repositories.districtRepository,
		s.repositories.shopCourierRepository,
		s.repositories.walletRepository,
		s.repositories.promotionRepository,
	)
	s.usecases.discoveryUsecase = usecase.NewDiscoveryUsecase(s.repositories.productRepository, s.repositories.reviewRepository)
	s.usecases.sellerPageUsecase = usecase.NewSellerPageUsecase(
		s.repositories.sellerPageRepository,
		s.repositories.reviewRepository,
	)
	s.usecases.wishlistUseCase = usecase.NewWishlistUsecase(s.repositories.wishlistRepository, s.repositories.productRepository, s.repositories.reviewRepository)
	s.usecases.reviewUsecase = usecase.NewReviewUsecase(s.repositories.reviewRepository, s.repositories.productRepository)
	s.usecases.promotionUsecase = usecase.NewPromotionRepository(s.repositories.promotionRepository, s.repositories.shopRepository)
}

func (s *server) initRESTHandler(logger dependency.Logger, config dependency.Config) {
	s.r = gin.Default()
	s.r.ContextWithFallback = true
	s.r.Use(
		middleware.CORS(config),
		middleware.ErrorHandler(),
		middleware.RequestID(),
		middleware.Logger(logger),
	)

	resthandler.NewAuthHandler(s.v, s.usecases.authUsecase, s.cfg).Route(s.r)
	resthandler.NewProfileHandler(s.usecases.accountAddressUsecase, s.cfg, s.v).Route(s.r)
	resthandler.NewHomePageHandler(s.v, s.usecases.homepageUsecase, s.cfg).Route(s.r)
	resthandler.NewShopHandler(s.usecases.shopUsecase, s.cfg, s.v).Route(s.r)
	resthandler.NewProductPageHandler(s.v, s.usecases.productPageUsecase, s.usecases.discoveryUsecase, s.cfg).Route(s.r)
	resthandler.NewDropdownHandler(s.v, s.usecases.dropdownUsecase, s.cfg).Route(s.r)
	resthandler.NewCartHandler(s.v, s.usecases.cartUsecase, s.cfg).Route(s.r)
	resthandler.NewWalletHandler(s.v, s.usecases.walletUsecase, config).Route(s.r)
	resthandler.NewOrderSellerHandler(s.v, s.usecases.orderSellerUsecase, config, s.usecases.orderUsecase).Route(s.r)
	resthandler.NewOrderHandler(s.usecases.orderUsecase, s.cfg, s.v).Route(s.r)
	resthandler.NewCheckoutHandler(s.v, s.usecases.checkoutUsecase, s.cfg).Route(s.r)
	resthandler.NewSellerPageHandler(s.usecases.sellerPageUsecase, s.cfg, s.v).Route(s.r)
	resthandler.NewWishlistHandler(s.usecases.wishlistUseCase, s.cfg, s.v).Route(s.r)
	resthandler.NewReviewHandler(s.usecases.reviewUsecase, s.cfg, s.v).Route(s.r)
	resthandler.NewShopPromotionHandler(s.usecases.promotionUsecase, s.cfg, s.v).Route(s.r)

	s.r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"message": "page not found"})
	})
}

func (s *server) startRESTServer(cfg dependency.Config) *http.Server {
	srv := http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Rest.Port),
		Handler: s.r,
	}

	go func() {
		log.Printf("REST server is running on port %d", cfg.Rest.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	return &srv
}

func initGracefulShutdown(restSrv *http.Server, cfg dependency.Config) {
	quit := make(chan os.Signal, 2)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.App.GracefulTimeout)*time.Second)
	defer cancel()

	// stop resthandler server
	if err := restSrv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	<-ctx.Done()
	log.Println("Server exiting")
}

func InitApp(db *sqlx.DB, rc *redis.Client, cfg dependency.Config, logger dependency.Logger) {
	s := server{
		v:   validator.New(),
		cfg: cfg,
	}

	shared.ValidatorUseJSONName(s.v)

	s.initRepository(db, rc, cfg)
	s.initUsecase(rc)
	s.initRESTHandler(logger, cfg)

	restSrv := s.startRESTServer(cfg)

	initGracefulShutdown(restSrv, cfg)
}
