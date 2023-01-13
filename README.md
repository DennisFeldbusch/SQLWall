# Overview

This is a Web Application Firewall (WAF) that aims to protect your web application from SQL injection attacks. It does this by checking the query part of the URL with regex to ensure it is valid.

# Features

- checks the query part of the URL with regular expressions
- configurable
- logging of blocked requests

# Installation

# Configuration

## Environment variables

| Variable      | Description                        | Default                 |
| ------------- | ---------------------------------- | ----------------------- |
| `PORT`        | Port to listen on                  | `8080`                  |
| `DESTINATION` | Destination to forward requests to | `http://127.0.0.1:8081` |

# Usage
