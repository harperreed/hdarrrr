# 📸 HDR Image Processor

## 📋 Summary of Project

Welcome to **HDR Image Processor**! This project is designed to create high dynamic range (HDR) images from three different exposure levels: low, mid, and high. Leveraging powerful image processing techniques, users can generate stunning HDR images with improved detail and color depth. 

With a user-friendly command-line interface, this tool is simple yet effective, making it accessible for developers and photography enthusiasts alike. Built with Go, this application is modular and extensible, allowing for the integration of different tone mapping algorithms.

## 🚀 How to Use

1. **Prerequisites**:
   - Ensure you have [Go](https://golang.org/dl/) installed on your system.
   - Prepare three images of the same dimensions but with varying exposure levels (low, mid, high).

2. **Clone the Repository**:
   ```bash
   git clone https://github.com/harperreed/hdarrrr.git
   cd hdarrrr
   ```

3. **Run the Application**:
   Use the following command to execute the program, replacing the paths with your image file paths:
   ```bash
   go run ./cmd/hdarrrr/main.go -low path/to/low_exposure.jpg -mid path/to/mid_exposure.jpg -high path/to/high_exposure.jpg -output path/to/output_image.jpg
   ```

   - The `-low`, `-mid`, and `-high` flags are required for specifying the input images.
   - The `-output` flag allows you to specify the name of the output HDR image. If omitted, the default will be `hdr_output.jpg`.

4. **Check Output**:
   After running the command, you should see a message indicating that the HDR image was successfully saved at the specified location!

## 🖼️ Using the Graphical Interface

1. **Install Fyne**:
   Ensure you have Fyne installed by running:
   ```bash
   go get fyne.io/fyne/v2
   ```

2. **Run the GUI Application**:
   Use the following command to execute the GUI application:
   ```bash
   go run ./cmd/hdarrrr/gui.go
   ```

3. **Select Images**:
   - Click on the "Select Image" buttons to choose three images with different exposure levels.
   - You can also drag and drop images directly into the application.

4. **Process Images**:
   - Click the "Process" button to start the HDR processing.
   - A progress bar will show the status of the image processing.
   - If there are any errors, user-friendly error messages will be displayed.

5. **View Result**:
   - Once the processing is complete, the resulting HDR image will be displayed.
   - You can zoom in to inspect image details closely.

## ⚙️ Tech Info

- **Language**: Go (Golang)
- **Dependencies**:
  - `github.com/mdouchement/hdr` for HDR processing.
  - `fyne.io/fyne/v2` for the graphical interface.
- **File Structure**:
  - `cmd/`: Contains the command line executable entry point and GUI implementation.
  - `pkg/`: Utilities and imaging packages.
  - `go.mod`: Module configuration file for managing dependencies.

### Features Included:
- Image loading and saving in PNG and JPEG formats.
- HDR image creation using `github.com/mdouchement/hdr`.
- Unit tests to ensure functionality across various components.
- Simple graphical interface using Fyne for image selection and processing.

### Contributing:
Contributions are welcome! If you have suggestions or improvements, feel free to create an issue or submit a pull request.

### License:
This project is open-source and available under the MIT License. See the `LICENSE` file for more details.

---

👨‍💻 Happy Coding! If you encounter any issues or have questions, please don't hesitate to reach out or submit an issue on the repository!
