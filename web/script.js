const backendUrl = "http://localhost:8080";

document.getElementById('submit-button').addEventListener('click', async function (event) {
    event.preventDefault();

    const longUrl = document.getElementById('longUrl').value;
    const loading = document.getElementById('loading');
    const button = document.getElementById('submit-button');
    const result = document.getElementById('result');
    const error = document.getElementById('error');
    const shortUrlElement = document.getElementById('shortUrl');

    button.disabled = true;  // Отключить кнопку на время запроса
    button.textContent = "Сокращение...";
    loading.style.display = "block"; // Показать спиннер

    try {
        const response = await fetch(`${backendUrl}/shorten`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ original: longUrl })
        });

        if (!response.ok) {
            throw new Error('Не удалось сократить ссылку');
        }

        const resultData = await response.json();
        shortUrlElement.textContent = resultData.short;
        result.classList.add('show');
        error.style.display = 'none';

    } catch (err) {
        error.textContent = err.message;
        error.style.display = 'block';
        result.classList.remove('show');
    } finally {
        button.disabled = false;
        button.textContent = "Сократить";
        loading.style.display = "none"; // Скрыть спиннер
    }
});

document.getElementById('shortUrl').addEventListener('click', function () {
    const shortUrl = document.getElementById('shortUrl').textContent;

    navigator.clipboard.writeText(shortUrl).then(function () {
        document.getElementById('copyNotice').style.display = 'block';

        setTimeout(() => {
            document.getElementById('copyNotice').style.display = 'none';
        }, 2000);
    }).catch(function (error) {
        console.error('Ошибка копирования: ', error);
    });
});
