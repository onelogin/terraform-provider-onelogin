# Getting Started with OneLogin Terraform Provider

This directory contains a minimal setup to get started with the OneLogin Terraform provider.

## Prerequisites

- [Terraform](https://www.terraform.io/downloads.html) installed (v0.13+)
- OneLogin account with API credentials
- API credentials with required permissions

## Setup

1. **Obtain OneLogin API Credentials**
   - Log in to your OneLogin admin portal
   - Navigate to Developers > API Credentials
   - Create credentials with appropriate permissions
   - Note the Client ID, Client Secret, and your subdomain

2. **Configure the Provider**
   - Open `main.tf`
   - Replace `YOUR_CLIENT_ID`, `YOUR_CLIENT_SECRET`, and `YOUR_SUBDOMAIN` with your actual credentials
   - The subdomain is the part before `.onelogin.com` in your OneLogin URL

3. **Initialize Terraform**
   ```
   terraform init
   ```

   > **Note:** If you encounter issues with version availability in the Terraform Registry, see the [Manual Installation](#manual-installation-from-github) section below.

4. **Apply Configuration**
   ```
   terraform apply
   ```

## Example Resources

The `main.tf` file includes a commented example of creating a OneLogin user. Uncomment it to use, or consult the [documentation](https://registry.terraform.io/providers/onelogin/onelogin/latest/docs) for more resource types.

## Manual Installation from GitHub

If you encounter issues with the provider availability in the Terraform Registry, you can manually install the provider from GitHub:

### Option 1: Download Pre-built Binaries

1. Create the plugins directory (adjust OS and architecture as needed):
   ```bash
   # For macOS Intel
   mkdir -p ~/.terraform.d/plugins/registry.terraform.io/onelogin/onelogin/0.6.0/darwin_amd64

   # For macOS Apple Silicon (M1/M2)
   mkdir -p ~/.terraform.d/plugins/registry.terraform.io/onelogin/onelogin/0.6.0/darwin_arm64

   # For Linux
   mkdir -p ~/.terraform.d/plugins/registry.terraform.io/onelogin/onelogin/0.6.0/linux_amd64

   # For Windows
   mkdir -p %APPDATA%\terraform.d\plugins\registry.terraform.io\onelogin\onelogin\0.6.0\windows_amd64
   ```

2. Download the provider binary from GitHub releases:
   ```bash
   # For macOS Intel
   curl -L https://github.com/onelogin/terraform-provider-onelogin/releases/download/v0.6.0/terraform-provider-onelogin_0.6.0_darwin_amd64.zip -o /tmp/provider.zip

   # For macOS Apple Silicon
   curl -L https://github.com/onelogin/terraform-provider-onelogin/releases/download/v0.6.0/terraform-provider-onelogin_0.6.0_darwin_arm64.zip -o /tmp/provider.zip

   # For Linux
   curl -L https://github.com/onelogin/terraform-provider-onelogin/releases/download/v0.6.0/terraform-provider-onelogin_0.6.0_linux_amd64.zip -o /tmp/provider.zip

   # For Windows (using PowerShell)
   Invoke-WebRequest -Uri https://github.com/onelogin/terraform-provider-onelogin/releases/download/v0.6.0/terraform-provider-onelogin_0.6.0_windows_amd64.zip -OutFile $env:TEMP\provider.zip
   ```

3. Extract and install:
   ```bash
   # For Unix-based systems (macOS/Linux)
   unzip /tmp/provider.zip -d /tmp/provider

   # Check the actual contents and structure of the extracted files
   ls -la /tmp/provider

   # For macOS Intel (adjust filenames based on what you see in the previous step)
   cp /tmp/provider/terraform-provider-onelogin_v0.6.0 ~/.terraform.d/plugins/registry.terraform.io/onelogin/onelogin/0.6.0/darwin_amd64/
   chmod +x ~/.terraform.d/plugins/registry.terraform.io/onelogin/onelogin/0.6.0/darwin_amd64/terraform-provider-onelogin_v0.6.0

   # For macOS Apple Silicon (adjust filenames based on what you see)
   cp /tmp/provider/terraform-provider-onelogin_v0.6.0 ~/.terraform.d/plugins/registry.terraform.io/onelogin/onelogin/0.6.0/darwin_arm64/
   chmod +x ~/.terraform.d/plugins/registry.terraform.io/onelogin/onelogin/0.6.0/darwin_arm64/terraform-provider-onelogin_v0.6.0

   # For Linux (adjust filenames based on what you see)
   cp /tmp/provider/terraform-provider-onelogin_v0.6.0 ~/.terraform.d/plugins/registry.terraform.io/onelogin/onelogin/0.6.0/linux_amd64/
   chmod +x ~/.terraform.d/plugins/registry.terraform.io/onelogin/onelogin/0.6.0/linux_amd64/terraform-provider-onelogin_v0.6.0

   # For Windows (using PowerShell) (adjust filenames based on what you see)
   Expand-Archive -Path $env:TEMP\provider.zip -DestinationPath $env:TEMP\provider
   Get-ChildItem -Path $env:TEMP\provider -Recurse
   Copy-Item $env:TEMP\provider\terraform-provider-onelogin_v0.6.0.exe $env:APPDATA\terraform.d\plugins\registry.terraform.io\onelogin\onelogin\0.6.0\windows_amd64\
   ```

   > **Note:** The exact structure and filenames inside the zip may vary between releases. After extraction, check the actual contents with `ls` or `Get-ChildItem` and adjust the copy commands accordingly. The binary may be in a subdirectory or have a slightly different name.

### Option 2: Build from Source

If you have the source code and Go installed:

1. Clone the repository:
   ```bash
   git clone https://github.com/onelogin/terraform-provider-onelogin.git
   cd terraform-provider-onelogin
   ```

2. Build and install:
   ```bash
   make sideload
   ```

3. This will compile the provider and install it in the correct location for Terraform to find it.

After manual installation, run `terraform init` in your project directory, and Terraform should find the locally installed provider.

## Security Note

Do not commit your API credentials to version control. Consider using environment variables or another secure method to manage secrets.

Example with environment variables:
```hcl
provider "onelogin" {
  client_id     = "${env("ONELOGIN_CLIENT_ID")}"
  client_secret = "${env("ONELOGIN_CLIENT_SECRET")}"
  subdomain     = "${env("ONELOGIN_SUBDOMAIN")}"
}
```