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
        const response = await fetch('/backend.cgi/save', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(data)
        });
        const result = await response.json();

        if (response.ok) {
            alert('Успешно отправлено: ' + result.message);
        } else {
            alert('Ошибка: ' + (result.error || 'Что-то пошло не так'));
        }
     } catch (err) {
            alert('Сервер недоступен');
            console.error(err);
        }
    }
)