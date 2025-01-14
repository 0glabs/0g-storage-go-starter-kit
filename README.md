# 0G Storage CLI

A command-line tool for uploading and downloading files using the 0G Storage network.

## Usage

### Upload a file
```bash
go run main.go -key YOUR_PRIVATE_KEY -upload path/to/file
```

### Download a file
```bash
go run main.go -key YOUR_PRIVATE_KEY -download ROOT_HASH -output path/to/save
```

### Options
- `-key`: Your private key for transactions (required)
- `-upload`: Path to the file you want to upload
- `-download`: Root hash of the file to download
- `-output`: Path where to save the downloaded file
- `-turbo`: Use Turbo endpoint for faster but more expensive operations (optional)

## Examples

Upload a file:
```bash
go run main.go -key 0x123... -upload ./myfile.txt
```

Download a file:
```bash
go run main.go -key 0x123... -download 0xabc... -output ./downloaded.txt
```

## Notes
- Private key is required for all operations
- Default replica count is set to 3
- Files are verified during download
- Turbo mode is available for faster operations at higher cost 