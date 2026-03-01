function eraseAll() {
    localStorage.clear();
    const form = document.querySelector("form");
    form.reset();
    Array.from(favlang.options).forEach(option => option.selected = false);
}

document.getElementById('form').addEventListener('submit', async (e) => {
    e.preventDefault();
    const formData = new FormData(e.target);
    const data = {
        surname: formData.get('surname'),
        name: formData.get('name'),
        midname: formData.get('midname'),
        phone: formData.get('phone'),
        email: formData.get("email"),
        bday: formData.get("bday"),
        gender: parseInt(formData.get("gender")),
        favlangs: formData.getAll('favlang[]').map(id => parseInt(id)),
        bio: formData.get('bio')
    };
    try {
        const response = await fetch('/backend_lab_3/backend.cgi/save', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(data)
        });
        const result = await response.json();

        if (response.ok) {
            alert('Успешно отправлено: ' + result.message);
            eraseAll();
        } else {
            alert('Ошибка: ' + (result.error || 'Что-то пошло не так'));
        }
    } catch (err) {
        alert('Сервер недоступен');
        console.error(err);
    }
}
)

const formContainer = document.querySelector("form");

formContainer.addEventListener("change", (e) => {
    const el = e.target;

    if (el.id === "favlangs") {
        const values = Array.from(el.selectedOptions).map(opt => opt.value);
        localStorage.setItem("favlangs", JSON.stringify(values));
    } else if (el.type === "radio") {
        localStorage.setItem("gender1", gender1.checked);
        localStorage.setItem("gender2", gender2.checked);
    } else {
        localStorage.setItem(el.id, el.value);
    }
});

const button = document.getElementById("submit-btn");
const checkbox = document.getElementById("check");
checkbox.addEventListener('change', (e) => {
    const isChecked = e.target.checked;
    if (isChecked && button.getAttribute("disabled") !== null) button.removeAttribute("disabled");
    else button.setAttribute("disabled", "");
});
window.addEventListener("DOMContentLoaded", () => {
    const form = document.querySelector("body");
    form.classList.remove("no-fouc");
    const textFields = ["surname", "name", "midname", "email", "bday", "phone", "bio"];
    textFields.forEach(id => {
        const val = localStorage.getItem(id);
        if (val) document.getElementById(id).value = val;
    });

    if (localStorage.getItem("gender1") === "true") gender1.checked = true;
    if (localStorage.getItem("gender2") === "true") gender2.checked = true;

    const savedLangs = localStorage.getItem("favlangs");
    if (savedLangs) {
        const langsArray = JSON.parse(savedLangs);
        Array.from(favlang.options).forEach(option => {
            if (langsArray.includes(option.value)) {
                option.selected = true;
            }
        });
    }
});
