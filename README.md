# BlubberBytes

**BlubberBytes** is a distributed file-sharing and HTTP proxy tool that allows users to share files by hash, set up a public gateway, and more.

## Prerequisites

Make sure you have the following installed on your system:

- **Node.js** (v14 or newer): Download from [nodejs.org](https://nodejs.org/).
- **npm** (Node Package Manager): Comes with Node.js.
- **Electron**: You do **not** need to install Electron globally; it will be installed automatically as part of the project dependencies.

## Getting Started

Follow these steps to set up and run the application locally:

### Step 1: Clone the Repository

First, clone the repository to your local machine:

```bash
git clone https://github.com/faefey/blubberbytes.git
```

### Step 2: Set Up the Client

Navigate to the `client` directory:

```bash
cd blubberbytes/client
```

### Step 3: Install Dependencies

Install the required dependencies by running:

```bash
npm install
```

> **Note**: This step will also install Electron as defined in the `package.json` file, so no separate installation is needed.

### Step 4: Build the Project

Once the dependencies are installed, build the project with:

```bash
npm run build
```

### Step 5: Run Electron

After building, you can run the Electron app:

```bash
npm run electron
```

This command will start the Electron desktop application.

## Running the Application

After completing the steps above, the GUI should open up automatically. 
