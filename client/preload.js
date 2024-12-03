const { contextBridge, webUtils } = require("electron")

contextBridge.exposeInMainWorld("electron", {
  pathForFile: (file) => webUtils.getPathForFile(file),
  posixPathForFile: (file) => {
    const path = webUtils.getPathForFile(file);
    // On Windows we get back a path with forward slashes.
    // Note that we don't have access to the path module in the preload script.
    return navigator.platform.toLowerCase().includes("win") ? path.split("\\").join("/") : path
  }
})