# HyperZoop Project

HyperZoop is a project focused on passwordless authentication. 

## Introduction

In the era of increasing security threats, the importance of secure authentication has never been more paramount. Passwordless authentication has emerged as a secure and user-friendly method for verifying users.

## What is Passwordless Authentication?

Passwordless authentication is a type of authentication where users do not need to log in with passwords. Instead of passwords, the system verifies the user by using something they have, such as a device, or something they are, like a fingerprint.

## How HyperZoop Implements Passwordless Authentication

In the HyperZoop project, we implement passwordless authentication by sending a unique link to the user's registered email. When the user clicks on this link, they are authenticated and logged into the system. This eliminates the need for users to remember complex passwords, and at the same time, enhances the security of the system.

Note: Please replace this section with the actual method your project uses for passwordless authentication.

## Getting Started

To get started with the HyperZoop project, follow these steps:

1. Clone the repository
2. Install the dependencies
3. Start the server

## Conclusion

HyperZoop project is a step forward in the realm of secure and user-friendly authentication. With passwordless authentication, we hope to provide users with an easy-to-use and secure method for accessing their accounts.


## .env variables

- env="dev" #describe application environment
- dev_db="" #can be replace for single db variable
- stagging_db="" #can be replace for single db variable
- prod_db="" #can be replace for single db variable
- redis_url="" #redis db for store magic_link/cache
- log_file="app.log" #store localhost logs into file
- app_host="localhost" 
- verify_host="http://localhost:3000"
- token_secret="" #jwt token