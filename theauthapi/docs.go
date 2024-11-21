// Package theauthapi provides a Go client for accessing the TheAuthAPI.
//
// Usage:
//
//	client := theauthapi.NewClient("YOUR_ACCESS_KEY")
//	
//	// Use the client to make API calls
//	isValid, err := client.ApiKeys.IsValidKey(ctx, "EXAMPLE_API_KEY")
package theauthapi