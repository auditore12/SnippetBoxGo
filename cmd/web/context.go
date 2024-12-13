// /In this file we’ll define a custom contextKey type and an isAuthenticatedContextKey variable, so that we have a unique key
// we can use to store and retrieve the authentication status from a request context (without the risk of naming collisions).
package main

type contextKey string
type currentUser string

const isAuthenticatedContextKey = contextKey("isAuthenticated")
const CurrentUserContexkey = contextKey("CurrentUser")
