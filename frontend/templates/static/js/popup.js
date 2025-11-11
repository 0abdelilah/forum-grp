function openPopup() {
    const overlay = document.getElementById("overlay");
    const popup = document.getElementById("popup");
    overlay.style.display = "flex";
    popup.style.display = "block";
}

function closePopup() {
    const overlay = document.getElementById("overlay");
    const popup = document.getElementById("popup");
    overlay.style.display = "none";
    popup.style.display = "none";
}
