export const detectKeyPress = (keyPressCallback) => {
  document.onkeyup = function (e) {
    if ((e.ctrlKey || e.metaKey) && (e.key === "1" || e.key === "2" || e.key === "3")) {
      console.log("Key pressed: ",e.key)
      keyPressCallback(e.key);
    }
  };
};