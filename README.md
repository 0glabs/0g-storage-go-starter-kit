# 0G Storage API Sandbox

This sandbox demonstrates how to use the 0G Storage API for uploading and downloading files.

## Getting Started in CodeSandbox

1. Click the "Fork" button to create your own sandbox
2. Create a new file named `.env` in the root directory:
   - Click the "+" button next to "Files"
   - Name it `.env`
   - Add your private key:
   ```
   PRIVATE_KEY=your_private_key_here
   ```

3. Open the terminal in CodeSandbox:
   - Click "View" -> "Terminal" or press `` Ctrl + ` ``
   - Run the server:
   ```bash
   go run main.go
   ```

4. Click the "Open in New Window" button in the browser preview to access the Swagger UI

## API Endpoints

The server provides two main endpoints:

### 1. Upload File
- **Endpoint**: `POST /api/v1/upload`
- **Input**: File (multipart/form-data)
- **Output**: Root hash and transaction hash
- **Try it**: Use the Swagger UI's "Try it out" button to upload a file

### 2. Download File
- **Endpoint**: `GET /download/{root_hash}`
- **Input**: Root hash from upload response
- **Output**: Original file
- **Try it**: Use the root hash from an upload to download the file

## Example Usage

### Using Swagger UI (Recommended)
1. Open the browser preview in a new window
2. Use the interactive Swagger UI to:
   - Upload files using the `/upload` endpoint
   - Download files using the `/download/{root_hash}` endpoint

### Using cURL (Advanced)
```bash
# Upload
curl -X POST http://localhost:8080/api/v1/upload \
  -F "file=@/path/to/your/file"

# Download
curl -O http://localhost:8080/api/v1/download/{root_hash}
```

## Notes
- Private key is required for uploading files
- Files are temporarily stored during transfer
- Default replica count is set to 1 for testing
- The server automatically redirects to Swagger UI for easy testing 