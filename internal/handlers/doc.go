// Package handlers implements the HTTP request processors for the application.
// It uses the Gin framework to handle routing, request binding, and response serialization.
//
// Each handler typically follows a pattern of:
// 1. Binding and validating the request payload using [dto] structures.
// 2. Invoking the appropriate [service] method to perform business logic.
// 3. Serializing the result or an error using the [RespondWithError] helper.
package handlers
