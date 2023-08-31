# Funds | Crowdfunding Platform

## About Funds

Fund is a crowdfunding platform that enables clients to post their projects in search of funding. With Fund, clients can create campaigns and showcase their ideas to a community of potential investors. This project built with Go

## Prerequisite

-   Go 1.18
-   PostgreSQL

## How to use this repository

-   Rename the `.env.example` file into `.env`.
-   Fill the value in `.env` file based on your case.
-   Run Http development server
    ```bash
    make dev
    ```

## Architecture

This repository implements Clean Architecture and a bit of SOLID principles.

Router -> Handler -> Service -> Repository

## Documentation

-   [Postman documentation](https://documenter.getpostman.com/view/16615700/2s93CNMCg8)
