// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"net/http"
	"strings"

	errors "github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"
	middleware "github.com/go-openapi/runtime/middleware"

	"github.com/lacchain/credential-server/restapi/operations"
	"github.com/lacchain/credential-server/restapi/operations/credential"
	"github.com/lacchain/credential-server/restapi/operations/did"

	"github.com/lacchain/credential-server/business"
	"github.com/rs/cors"
)

//go:generate swagger generate server --target ../../CredentialProvider --name CredentialProvider --spec ../swagger/swaggerui/swagger.json

func configureFlags(api *operations.CredentialProviderAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.CredentialProviderAPI, config *ServerConfig) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	api.XMLProducer = runtime.XMLProducer()

	api.CredentialCreateCredentialHandler = credential.CreateCredentialHandlerFunc(func(params credential.CreateCredentialParams) middleware.Responder {
		response, err := business.CreateCredential(params.Body, config.Node, config.Issuer, config.PrivateKey, config.Proof.Verification, config.Repository.Address)
		if err != nil {
			return credential.NewCreateCredentialBadRequest()
		}
		return credential.NewCreateCredentialOK().WithPayload(response)
	})
	api.CredentialRenewalCredentialHandler = credential.RenewalCredentialHandlerFunc(func(params credential.RenewalCredentialParams) middleware.Responder {
		response, err := business.RenewCredential(params.CredentialID, params.Body, config.Node, config.Issuer, config.PrivateKey, config.Proof.Verification, config.Repository.Address)
		if err != nil {
			return credential.NewRenewalCredentialBadRequest()
		}
		return credential.NewRenewalCredentialOK().WithPayload(response)
	})
	api.DidValidateDidHandler = did.ValidateDidHandlerFunc(func(params did.ValidateDidParams) middleware.Responder {
		validate, err := business.ValidateDid(params.Did)
		if !validate || err != nil {
			return did.NewValidateDidBadRequest()
		}
		return did.NewValidateDidOK()
	})
	api.CredentialVerifyCredentialHandler = credential.VerifyCredentialHandlerFunc(func(params credential.VerifyCredentialParams) middleware.Responder {
		response, err := business.VerifyCredential(params.Body, config.Node, config.Address, config.Proof.Verification)
		if err != nil {
			return credential.NewVerifyCredentialBadRequest()
		}
		return credential.NewVerifyCredentialOK().WithPayload(response)
	})
	api.CredentialRevokeCredentialHandler = credential.RevokeCredentialHandlerFunc(func(params credential.RevokeCredentialParams) middleware.Responder {
		err := business.RevokeCredential(params.CredentialID, config.Node, config.PrivateKey, config.Proof.Verification, config.Repository.Address)
		if err != nil {
			return credential.NewRevokeCredentialBadRequest()
		}
		return credential.NewRevokeCredentialOK()
	})

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix"
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {

	corsHandler := cors.New(cors.Options{
		Debug:          false,
		AllowedHeaders: []string{"*"},
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{},
		MaxAge:         1000,
	})
	return corsHandler.Handler(uiMiddleware(handler))
}

func uiMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Shortcut helpers for swagger-ui
		if r.URL.Path == "/swagger-ui" || r.URL.Path == "/api/help" {
			http.Redirect(w, r, "/swagger-ui/", http.StatusFound)
			return
		}
		// Serving ./swagger-ui/
		if strings.Index(r.URL.Path, "/swagger-ui/") == 0 {
			http.StripPrefix("/swagger-ui/", http.FileServer(http.Dir("./swaggerui/"))).ServeHTTP(w, r)
			return
		}

		handler.ServeHTTP(w, r)
	})
}
