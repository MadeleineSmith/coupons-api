# Coupons API

## Getting started
1. Create a local postgres database called `coupons`
2. Fill out your database credentials in `example_config.json` and rename file to `config.json`
3. `cd` to each of `handlers`, `model/coupon`, `dbservices` and execute `ginkgo` to ensure that all unit tests are green
4. Run the application with `go build` followed by `./coupons` 