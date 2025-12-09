# V. Conclusion and Future Work

## Conclusion
This project successfully demonstrates the implementation of a ultimate secure and optimized performance steganography tool by combining **Least Significant Bit (LSB)** encoding with **AES-256-GCM** encryption with **Hashed Pass(SHA-256)**.
The tool effectively solves the problem of covert communication by allowing users to hide data in plain sight, leveraging the ubiquity of digital images.

## Future Work
While the current system is functional and secure, several enhancements could be made in future iterations:

### 1. Support for More Image Formats
Currently, the tool relies on PNG (lossless). Supporting **JPEG** would require more complex algorithms.

### 2. Advanced Key Derivation
Currently, SHA-256 is used to hash the password. To better resist brute-force attacks, a memory-hard function like **Argon2id** or **scrypt** should be implemented.

### 3. Adaptive Steganography
Instead of sequential LSB replacement (which can be statistically detected), the tool could implement **adaptive algorithms** that hide data in "noisy" areas of the image (edges, textures) where it is harder to detect.

### 4. Mobile Application
Developing a native mobile app (React Native or Flutter) would make the tool more accessible for secure communication on the go.

