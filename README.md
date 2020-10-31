# Description
Script that connects with an API that provides debt, payment_plan, and payment information. It enriches the data to provide the debt_id, debt_amount, next_payment_date, amount_owed

# Environment variables
```TRUEACCORD_API_URL```

## Example environment variables:
```bash
TRUEACCORD_API_URL=http://my-json-server.typicode.com/pink-cupcakes/TrueAccord
```

# To run the true_accord service
Requires go 1.13
```bash
cd TrueAccord
go run true_accord
```
Note: the executible binary is included and can be run directly. If it fails - check if the environment variables were set.