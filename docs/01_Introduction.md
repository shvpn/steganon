# I. Introduction / Background

## Project Goal
The primary goal of this project is to develop a secure and user-friendly **Steganography Tool** that allows users to hide secret text messages within digital images. Unlike traditional encryption, which scrambles data but leaves the communication visible, this tool aims to provide **covert communication** by concealing the very existence of the message with ultimate security in 2026.

## Problem Statement
In an era of pervasive digital surveillance, protecting the privacy of communication is increasingly challenging.
1.  **Visibility of Encryption**: Standard encrypted messages (like PGP or encrypted emails) clearly signal that sensitive information is being transmitted, potentially attracting unwanted attention from adversaries.
2.  **Complexity**: Many steganographic tools are command-line based or require complex installation, making them inaccessible to average users.
3.  **Lack of Security**: Basic steganography tools often hide plain text. If the hidden message is detected, it can be easily read.

## Proposed Solution
This project implements a web-based application that combines **Cryptography** and **Steganography** to address these issues:
-   **Obscurity**: Uses **Least Significant Bit (LSB)** steganography to embed data into the noise of an image, making it invisible to the human eye.
-   **Confidentiality**: Incorporates **AES-256-GCM** encryption to ensure that even if the hidden data is extracted, it cannot be read without the correct password.
-   **Accessibility**: Provides a modern, responsive web interface that works on any device with a browser.

## Motivation
The motivation behind this project is to explore the intersection of security and privacy. It serves as a practical application of cryptographic primitives (AES, SHA-256) and image processing algorithms, demonstrating how different security layers can be layered to provide robust protection for sensitive data.

## Related Cryptographic Concepts

### 1. Steganography
Derived from Greek words *steganos* (covered) and *graphein* (writing). It is the practice of concealing a file, message, image, or video within another file, message, image, or video.
-   **Carrier**: The signal, stream, or data file that hides the payload (in this case, a PNG image).
-   **Payload**: The information to be hidden (the secret message).

### 2. Least Significant Bit (LSB)
A common steganographic technique where the last bit of each pixel's byte is replaced with a bit of the secret message. Since the change is minimal (changing a color value by 1 out of 255), the alteration is imperceptible to the human visual system.

### 3. Advanced Encryption Standard (AES)
A symmetric block cipher chosen by the U.S. government to protect classified information. This project uses **AES-256**, which uses 256-bit keys.

### 4. Galois/Counter Mode (GCM)
A mode of operation for symmetric-key cryptographic block ciphers. It provides both **data authenticity (integrity)** and **confidentiality**. It ensures that the message has not been tampered with during transmission.

### 5. Hashing (SHA-256)
A cryptographic hash function that outputs a value that is 256 bits long. In this project, it is used for **Key Derivation**, converting a user-provided password of any length into a fixed-size 32-byte key required for AES-256.
