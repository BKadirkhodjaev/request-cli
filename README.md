# Request CLI

## Purpose

- Command line utility that performs parallel HTTP requests using a provided CSV file

## Command

- Create a CSV file in `import/data.csv` with a similar column structure as shown below

```csv
"po_id","vendor_id","po_number"
ceb20319-4029-413a-9acf-35b9a943bb07,a31a6f0f-6899-477d-9ec7-98df2531858e,"10000"
"1f315c3c-c49d-444a-8b61-1b21af8d9358",a31a6f0f-6899-477d-9ec7-98df2531858e,"10001"
"07beaba7-c367-4ea4-a23b-8a184376228a",a31a6f0f-6899-477d-9ec7-98df2531858e,"10002"
```  

- Build & Run

```shell
mkdir -p bin
env GOOS=windows GOARCH=amd64 go build -o ./bin . && ./bin/request-cli.exe 
```
