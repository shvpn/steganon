# IV. Usage Guide

## User Guide

### Feature 1: Encoding (Hiding a Message)

1.  **Select Tab**: Click on the **"Encode"** tab.
2.  **Upload Image**: Click the file input area and select a **PNG** image.
    *   *Note: PNG is recommended because it is lossless. JPEG compression can destroy the hidden LSB data.*
3.  **Enter Message**: Type the secret message you want to hide in the text area.
4.  **Set Password (Optional)**:
    -   Enter a password to encrypt the message.
    -   If left blank, the message will be hidden but not encrypted.
5.  **Process**: Click the **"Encode Message"** button.
6.  **Download**: Once finished, a "Download Encoded Image" button will appear. Click it to save the result.

### Feature 2: Decoding (Reading a Message)

1.  **Select Tab**: Click on the **"Decode"** tab.
2.  **Upload Image**: Select the image that contains the hidden message.
3.  **Enter Password**:
    -   If the message was encrypted, you **must** enter the same password used during encoding.
    -   If no password was used, leave this field blank.
4.  **Process**: Click the **"Decode Message"** button.
5.  **View Result**: The hidden message will appear on the screen. You can use the "Copy" button to copy it to your clipboard.

## Example Scenario

1.  **Alice** wants to send a secret bank account number to **Bob**.
2.  Alice opens the tool, uploads a picture of a cat (`cat.png`).
3.  Alice types the account number and sets the password to `Secret123!`.
4.  Alice downloads the result (`encoded.png`) and emails it to Bob.
5.  **Bob** receives the email. To anyone else, it looks like just a cat picture.
6.  Bob opens the tool, uploads `encoded.png`, enters `Secret123!`, and retrieves the account number.
