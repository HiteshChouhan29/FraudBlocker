document.addEventListener('DOMContentLoaded', function() {
    const checkButton = document.getElementById('check-button');
    const checkPhone = document.getElementById('check-phone');
    const checkResult = document.getElementById('check-result');

    checkButton.addEventListener('click', function() {
        const phoneNumber = checkPhone.value.trim();
        if (!phoneNumber) {
            alert('Please enter a phone number');
            return;
        }
        fetch('/check?phone=' + encodeURIComponent(phoneNumber))
            .then(response => response.json())
            .then(data => {
                checkResult.classList.remove('hidden', 'blocked', 'not-blocked');
                if (data.blocked) {
                    alert('WARNING: This phone number is BLOCKED!');
                    checkResult.textContent = 'WARNING: This phone number is BLOCKED!';
                    checkResult.classList.add('blocked');
                } else {
                    checkResult.textContent = 'This phone number is not blocked.';
                    checkResult.classList.add('not-blocked');
                }
            })
            .catch(error => {
                console.error('Error:', error);
                alert('Error checking phone number');
            });
    });

    const addForm = document.querySelector('form[action="/add"]');
    if (addForm) {
        addForm.addEventListener('submit', function(event) {
            event.preventDefault();
            const phoneInput = addForm.querySelector('input[name="phone"]');
            const phoneNumber = phoneInput.value.trim();
            if (!phoneNumber) {
                alert('Please enter a phone number');
                return;
            }
            const phoneRegex = /^[+]?[\d\s\-()]{7,20}$/;
            if (!phoneRegex.test(phoneNumber)) {
                alert('Invalid phone number format');
                return;
            }
            const formData = new FormData(addForm);
            fetch('/add', {
                method: 'POST',
                body: formData
            })
                .then(response => {
                    if (response.status === 409) {
                        alert('This phone number is already blocked.');
                        throw new Error('Duplicate blocked number');
                    } else if (!response.ok) {
                        alert('Error adding phone number.');
                        throw new Error('Error adding phone number');
                    } else {
                        window.location.href = '/';
                    }
                })
                .catch(error => {
                    console.error('Add number error:', error);
                });
        });
    }

    const removeForms = document.querySelectorAll('form[action="/remove"]');
    removeForms.forEach(form => {
        form.addEventListener('submit', function(event) {
            event.preventDefault();
            const formData = new FormData(form);
            fetch('/remove', {
                method: 'POST',
                body: formData
            })
                .then(response => {
                    if (!response.ok) {
                        alert('Error removing phone number.');
                        throw new Error('Error removing phone number');
                    } else {
                        alert('Phone number removed successfully.');
                        window.location.reload();
                    }
                })
                .catch(error => {
                    console.error('Remove number error:', error);
                });
        });
    });
});