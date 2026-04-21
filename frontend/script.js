const API_BASE = '/api';

let currentUser = null;
let movies = [];
let screenings = [];
let reservations = [];
let selectedMovie = null;

function getAuthHeader() {
    const token = localStorage.getItem('token');
    return token ? { 'Authorization': `Bearer ${token}` } : {};
}

async function apiFetch(url, options = {}) {
    const response = await fetch(url, options);
    if (response.status === 401) {
        handleLogout();
        throw new Error('Session expired');
    }
    return response;
}

document.addEventListener('DOMContentLoaded', () => {
    const storedUser = localStorage.getItem('user');
    if (storedUser) {
        currentUser = JSON.parse(storedUser);
        showMainPage();
        loadMovies();
    } else {
        showAuthPage();
    }
});

        document.getElementById('loginForm').reset();
function renderMovies() {
    const moviesList = document.getElementById('moviesList');

    if (movies.length === 0) {
        moviesList.innerHTML = '<p>No movies available.</p>';
        return;
    }

    moviesList.innerHTML = movies.map(movie => `
        <div class="movie-card">
            <div class="movie-poster">
                ${movie.poster_url ? `<img src="${movie.poster_url}" alt="${escapeHtml(movie.title)} poster" />` : '<div class="poster-placeholder">No poster</div>'}
            </div>
            <div class="movie-header">
                <h3>${escapeHtml(movie.title)}</h3>
                <span class="movie-genre">${escapeHtml(movie.genre)}</span>
            </div>
            <p class="movie-description">${escapeHtml(movie.description)}</p>
            <div class="movie-meta">
                <span>Duration: ${movie.duration} min</span>
            </div>
            <div class="movie-actions">
                <button class="btn btn-primary" onclick="openBookingModal(${movie.id})">Book</button>
                    ${currentUser && currentUser.is_admin ? '' : ''}
            </div>
        </div>
    `).join('');
}
        showMainPage();
        loadMovies();
        loadReservations();
    } catch (error) {
        showAuthMessage('Network error', 'error');
        console.error('Login error:', error);
    }
}

async function handleRegister(e) {
    e.preventDefault();
    const fullName = document.getElementById('registerName').value;
    const email = document.getElementById('registerEmail').value;
    const password = document.getElementById('registerPassword').value;
    const isAdmin = document.getElementById('registerIsAdmin').checked;

    try {
        const response = await fetch(`${API_BASE}/register`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ email, password, full_name: fullName, is_admin: isAdmin })
        });

        if (!response.ok) {
            showAuthMessage('Registration failed.', 'error');
            return;
        }

        showAuthMessage('Registration successful! You can now log in.', 'success');
        document.getElementById('registerForm').reset();
        switchTab('login');
    } catch (error) {
        showAuthMessage('Network error', 'error');
        console.error('Register error:', error);
    }
}

function handleLogout() {
    localStorage.removeItem('token');
    localStorage.removeItem('user');
    currentUser = null;
    movies = [];
    screenings = [];
    reservations = [];
    showAuthPage();
}

function switchTab(tab) {
    const loginForm = document.getElementById('loginForm');
    const registerForm = document.getElementById('registerForm');
    const tabBtns = document.querySelectorAll('.tab-btn');

    if (tab === 'login') {
        loginForm.classList.add('active');
        registerForm.classList.remove('active');
        tabBtns[0].classList.add('active');
        tabBtns[1].classList.remove('active');
    } else {
        loginForm.classList.remove('active');
        registerForm.classList.add('active');
        tabBtns[0].classList.remove('active');
        tabBtns[1].classList.add('active');
    }
}

function showAuthPage() {
    document.getElementById('authPage').style.display = 'block';
    document.getElementById('mainPage').style.display = 'none';
}

function showMainPage() {
    document.getElementById('authPage').style.display = 'none';
    document.getElementById('mainPage').style.display = 'block';

    if (currentUser) {
        document.getElementById('userName').textContent = currentUser.full_name;
    authMessage.classList.add('show');
function renderAllScreenings() {
    const screeningsList = document.getElementById('screeningsList');
    if (!screenings.length) {
        screeningsList.innerHTML = '<p>No screenings available.</p>';
        return;
    }

    screeningsList.innerHTML = screenings.map(screening => {
        const date = new Date(screening.starts_at).toLocaleString('en-GB', { dateStyle: 'short', timeStyle: 'short' });
        const movie = movies.find(m => m.id === screening.movie_id);
        const movieTitle = movie ? escapeHtml(movie.title) : `Movie ${screening.movie_id}`;

        return `
            <div class="screening-card">
                <div class="screening-info">
                    <h3>${movieTitle}</h3>
                    <p>Screening: ${date}</p>
                    <p>Available seats: ${screening.available_seats} | Price: ${screening.price}</p>
                </div>
                            <div class="screening-actions">
                                ${currentUser && currentUser.is_admin ? `<button class="btn btn-danger" onclick="handleDeleteScreening(${screening.id})">Delete</button>` : ''}
                </div>
            </div>
        `;
    }).join('');
}

    setTimeout(() => {
        authMessage.classList.remove('show');
    }, 5000);
}

async function loadMovies() {
    try {
        const response = await fetch(`${API_BASE}/movies`, {
            cache: 'no-cache',
            headers: {
                'Cache-Control': 'no-cache'
            }
        });
        if (!response.ok) throw new Error('Failed to load movies');

        movies = await response.json();
        renderMovies();
        renderMovieOptionsForScreening();
    } catch (error) {
        console.error('Load movies error:', error);
    }
}

function renderMovieOptionsForScreening() {
    const movieSelect = document.getElementById('screeningMovieSelect');
    if (!movieSelect) return;

    if (!movies.length) {
        movieSelect.innerHTML = '<option value="">Add at least one movie first</option>';
        return;
    }

    movieSelect.innerHTML = movies.map(movie => `<option value="${movie.id}">${escapeHtml(movie.title)}</option>`).join('');
}

async function handleCreateMovie(e) {
    e.preventDefault();

    const title = document.getElementById('movieTitle').value.trim();
    const description = document.getElementById('movieDescription').value.trim();
    const duration = parseInt(document.getElementById('movieDuration').value, 10);
    const genre = document.getElementById('movieGenre').value.trim();
    const posterUrl = document.getElementById('moviePosterUrl').value.trim();
    const messageBox = document.getElementById('movieFormMessage');

    messageBox.textContent = '';
    messageBox.className = 'alert';

    try {
        const response = await apiFetch(`${API_BASE}/movies`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json', ...getAuthHeader() },
            body: JSON.stringify({ title, description, duration, genre, poster_url: posterUrl })
        });

        if (!response.ok) {
            const errorText = await response.text();
            messageBox.textContent = errorText || 'Failed to create movie.';
            messageBox.classList.add('error');
            return;
        }

        const newMovie = await response.json();
        movies.push(newMovie);
        renderMovies();
        renderMovieOptionsForScreening();

        document.getElementById('createMovieForm').reset();
        messageBox.textContent = 'Movie added successfully!';
        messageBox.classList.add('success');
        await loadMovies();
    } catch (error) {
        messageBox.textContent = 'Network error.';
        messageBox.classList.add('error');
        console.error('Create movie error:', error);
    }
}

async function handleDeleteMovie(movieId) {
    if (!confirm('Are you sure you want to delete this movie?')) return;

    try {
        const response = await apiFetch(`${API_BASE}/movies/${movieId}`, {
            method: 'DELETE',
            headers: { ...getAuthHeader() }
        });

        if (!response.ok) {
            const errorText = await response.text();
            alert(errorText || 'Failed to delete movie.');
            return;
        }

        movies = movies.filter(movie => movie.id !== movieId);
        renderMovies();
        renderMovieOptionsForScreening();
    } catch (error) {
        console.error('Delete movie error:', error);
        alert('A network error occurred while deleting the movie.');
    }
}

async function handleCreateScreening(e) {
    e.preventDefault();

    const movieId = parseInt(document.getElementById('screeningMovieSelect').value, 10);
    const startsAt = document.getElementById('screeningStartsAt').value;
    const availableSeats = parseInt(document.getElementById('screeningSeats').value, 10);
    const price = parseFloat(document.getElementById('screeningPrice').value);
    const messageBox = document.getElementById('screeningFormMessage');

    messageBox.textContent = '';
    messageBox.className = 'alert';

    if (!movieId || !startsAt || availableSeats <= 0 || price <= 0) {
        messageBox.textContent = 'Please fill in all fields.';
        messageBox.classList.add('error');
        return;
    }

    const isoStartsAt = new Date(startsAt).toISOString();

    try {
        const response = await apiFetch(`${API_BASE}/screenings`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json', ...getAuthHeader() },
            body: JSON.stringify({ movie_id: movieId, starts_at: isoStartsAt, available_seats: availableSeats, price })
        });

        if (!response.ok) {
            const errorText = await response.text();
            messageBox.textContent = errorText || 'Failed to create screening.';
            messageBox.classList.add('error');
            return;
        }

        document.getElementById('createScreeningForm').reset();
        messageBox.textContent = 'Screening added successfully!';
        messageBox.classList.add('success');
        await loadAllScreenings();
    } catch (error) {
        messageBox.textContent = 'Network error.';
        messageBox.classList.add('error');
        console.error('Create screening error:', error);
    }
}

function renderMovies() {
    const moviesList = document.getElementById('moviesList');

    if (movies.length === 0) {
        moviesList.innerHTML = '<p>No movies available.</p>';
        return;
    }

    moviesList.innerHTML = movies.map(movie => `
        <div class="movie-card">
            <div class="movie-poster">
                ${movie.poster_url ? `<img src="${movie.poster_url}" alt="${escapeHtml(movie.title)} poster" />` : '<div class="poster-placeholder">No poster</div>'}
            </div>
            <div class="movie-header">
                <h3>${escapeHtml(movie.title)}</h3>
                <span class="movie-genre">${escapeHtml(movie.genre)}</span>
            </div>
            <p class="movie-description">${escapeHtml(movie.description)}</p>
            <div class="movie-meta">
                <span>Duration: ${movie.duration} min</span>
            </div>
            <div class="movie-actions">
                <button class="btn btn-primary" onclick="openBookingModal(${movie.id})">Book</button>
            </div>
        </div>
    `).join('');
}

                ${currentUser && currentUser.is_admin ? `<button class="btn btn-danger" onclick="handleDeleteMovie(${movie.id})">Delete</button>` : ''}
async function openBookingModal(movieId) {
    selectedMovie = movies.find(movie => movie.id === movieId);
    if (!selectedMovie) return;

    document.getElementById('bookingModalTitle').textContent = `Book: ${selectedMovie.title}`;
    document.getElementById('bookingError').textContent = '';
    document.getElementById('seatCount').value = 1;
    document.getElementById('screeningSelect').value = '';

    try {
        const response = await fetch(`${API_BASE}/screenings?movie_id=${movieId}`);
        if (!response.ok) throw new Error('Failed to load screenings');

        screenings = await response.json();
        renderScreeningOptions();
        
        if (screenings.length > 0) {
            document.getElementById('screeningSelect').value = screenings[0].id;
        }
        
        document.getElementById('bookingModal').classList.add('show');
    } catch (error) {
        document.getElementById('bookingError').textContent = 'Failed to load screenings.';
        console.error('Load screenings error:', error);
    }
}

function renderScreeningOptions() {
    const screeningSelect = document.getElementById('screeningSelect');
    if (!screenings.length) {
        screeningSelect.innerHTML = '<option value="">No screenings available</option>';
        return;
    }

    screeningSelect.innerHTML = screenings.map(screening => {
        const date = new Date(screening.starts_at).toLocaleString('en-GB', { dateStyle: 'short', timeStyle: 'short' });
        return `<option value="${screening.id}">${date} - Available: ${screening.available_seats} - Price: ${screening.price}</option>`;
    }).join('');
}

function closeBookingModal() {
    document.getElementById('bookingModal').classList.remove('show');
}

async function handleCreateReservation(e) {
    e.preventDefault();

    const screeningId = parseInt(document.getElementById('screeningSelect').value, 10);
    const seats = parseInt(document.getElementById('seatCount').value, 10);

    if (!screeningId || seats <= 0) {
        document.getElementById('bookingError').textContent = 'Please select a screening and number of seats.';
        return;
    }

    try {
        const response = await apiFetch(`${API_BASE}/reservations`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                ...getAuthHeader()
            },
            body: JSON.stringify({ screening_id: screeningId, seats })
        });

        if (!response.ok) {
            const message = await response.text();
            document.getElementById('bookingError').textContent = message || 'Booking failed.';
            return;
        }

        closeBookingModal();
        loadReservations();
    } catch (error) {
        document.getElementById('bookingError').textContent = 'Network error.';
        console.error('Booking error:', error);
    }
}

async function loadReservations() {
    try {
        const response = await apiFetch(`${API_BASE}/reservations`, {
            headers: { ...getAuthHeader() }
        });
        if (!response.ok) throw new Error('Failed to load reservations');

        reservations = await response.json();
        renderReservations();
    } catch (error) {
        console.error('Load reservations error:', error);
    }
}

function renderReservations() {
    const reservationsList = document.getElementById('reservationsList');
    if (!reservations.length) {
        reservationsList.innerHTML = '<p>You have no reservations yet.</p>';
        return;
    }

    reservationsList.innerHTML = reservations.map(reservation => `
        <div class="reservation-card">
            <div>
                <h3>${escapeHtml(reservation.movie_title)}</h3>
                <p>${new Date(reservation.starts_at).toLocaleString('en-GB', { dateStyle: 'short', timeStyle: 'short' })}</p>
            </div>
            <div>
                <span>${reservation.seats} seat(s)</span>
                <span>${reservation.price}</span>
            </div>
        </div>
    `).join('');
}

function escapeHtml(text) {
    if (!text) return '';
    const map = {
        '&': '&amp;',
        '<': '&lt;',
        '>': '&gt;',
        '"': '&quot;',
        "'": '&#039;'
    };
    return String(text).replace(/[&<>"']/g, m => map[m]);
}

async function loadAllScreenings() {
    try {
        const response = await fetch(`${API_BASE}/screenings`, {
            cache: 'no-cache',
            headers: {
                'Cache-Control': 'no-cache'
            }
        });
        if (!response.ok) throw new Error('Failed to load screenings');

        screenings = await response.json();
        renderAllScreenings();
    } catch (error) {
        console.error('Load all screenings error:', error);
    }
}

function renderAllScreenings() {
    const screeningsList = document.getElementById('screeningsList');
    if (!screenings.length) {
        screeningsList.innerHTML = '<p>No screenings available.</p>';
        return;
    }

    screeningsList.innerHTML = screenings.map(screening => {
        const date = new Date(screening.starts_at).toLocaleString('en-GB', { dateStyle: 'short', timeStyle: 'short' });
        const movie = movies.find(m => m.id === screening.movie_id);
        const movieTitle = movie ? escapeHtml(movie.title) : `Movie ${screening.movie_id}`;

        return `
            <div class="screening-card">
                <div class="screening-info">
                    <h3>${movieTitle}</h3>
                    <p>Screening: ${date}</p>
                    <p>Available seats: ${screening.available_seats} | Price: ${screening.price}</p>
                </div>
                    <div class="screening-actions">
                        ${currentUser && currentUser.is_admin ? `<button class="btn btn-danger" onclick="handleDeleteScreening(${screening.id})">Delete</button>` : ''}
                    </div>
            </div>
        `;
    }).join('');
}

                    ${currentUser && currentUser.is_admin ? `<button class="btn btn-danger" onclick="handleDeleteScreening(${screening.id})">Delete</button>` : ''}
async function handleDeleteScreening(screeningId) {
    if (!confirm('Are you sure you want to delete this screening?')) return;

    try {
        const response = await apiFetch(`${API_BASE}/screenings/${screeningId}`, {
            method: 'DELETE',
            headers: { ...getAuthHeader() }
        });

        if (!response.ok) {
            const errorText = await response.text();
            alert(errorText || 'Failed to delete screening.');
            return;
        }

        screenings = screenings.filter(s => s.id !== screeningId);
        renderAllScreenings();
    } catch (error) {
        console.error('Delete screening error:', error);
        alert('A network error occurred while deleting the screening.');
    }
}

async function loadAllReservations() {
    try {
        const response = await apiFetch(`${API_BASE}/reservations`, {
            headers: { ...getAuthHeader() }
        });
        if (!response.ok) throw new Error('Failed to load reservations');

        reservations = await response.json();
        renderAllReservations();
    } catch (error) {
        console.error('Load all reservations error:', error);
    }
}

function renderAllReservations() {
    const reservationsList = document.getElementById('reservationsList');
    if (!reservations.length) {
        reservationsList.innerHTML = '<p>No reservations found.</p>';
        return;
    }

    reservationsList.innerHTML = reservations.map(reservation => `
        <div class="reservation-card">
            <div>
                <h3>${escapeHtml(reservation.movie_title)}</h3>
                <p>${new Date(reservation.starts_at).toLocaleString('en-GB', { dateStyle: 'short', timeStyle: 'short' })}</p>
                <p class="text-light">Seats: ${reservation.seats}</p>
            </div>
            <div>
                <span>${reservation.price}</span>
                <button class="btn btn-danger btn-small" onclick="handleDeleteReservation(${reservation.id})">Delete</button>
            </div>
        </div>
    `).join('');
}

async function handleDeleteReservation(reservationId) {
    if (!confirm('Are you sure you want to delete this reservation?')) return;

    try {
        const response = await apiFetch(`${API_BASE}/reservations/${reservationId}`, {
            method: 'DELETE',
            headers: { ...getAuthHeader() }
        });

        if (!response.ok) {
            const errorText = await response.text();
            alert(errorText || 'Failed to delete reservation.');
            return;
        }

        reservations = reservations.filter(r => r.id !== reservationId);
        renderAllReservations();
    } catch (error) {
        console.error('Delete reservation error:', error);
        alert('A network error occurred while deleting the reservation.');
    }
}

function showNewReservationSection() {
    const section = document.getElementById('newReservationSection');
    section.style.display = section.style.display === 'none' ? 'block' : 'none';
    if (section.style.display === 'block') {
        renderMoviesForNewReservation();
    }
}

async function renderMoviesForNewReservation() {
    try {
        const response = await fetch(`${API_BASE}/movies`);
        if (!response.ok) throw new Error('Failed to load movies');
        
        const moviesData = await response.json();
        const container = document.getElementById('newReservationMoviesList');
        
        if (!moviesData || moviesData.length === 0) {
            container.innerHTML = '<p>No movies available.</p>';
            return;
        }
        
        container.innerHTML = moviesData.map(movie => `
            <div class="movie-card">
                <div class="movie-poster">
                    ${movie.poster_url ? `<img src="${movie.poster_url}" alt="${escapeHtml(movie.title)} poster" />` : '<div class="poster-placeholder">No poster</div>'}
                </div>
                <div class="movie-header">
                    <h3>${escapeHtml(movie.title)}</h3>
                    <span class="movie-genre">${escapeHtml(movie.genre)}</span>
                </div>
                <p class="movie-description">${escapeHtml(movie.description)}</p>
                <div class="movie-meta">
                    <span>Duration: ${movie.duration} min</span>
                </div>
                <div class="movie-actions">
                    <button class="btn btn-primary" onclick="openNewReservationModal(${movie.id})">Book</button>
                </div>
            </div>
        `).join('');
    } catch (error) {
        console.error('Load movies error:', error);
    }
}

async function openNewReservationModal(movieId) {
    selectedMovie = movies.find(movie => movie.id === movieId);
    if (!selectedMovie) return;

    document.getElementById('bookingModalTitle').textContent = `Book: ${selectedMovie.title}`;
    document.getElementById('bookingError').textContent = '';
    document.getElementById('seatCount').value = 1;
    document.getElementById('screeningSelect').value = '';

    try {
        const response = await fetch(`${API_BASE}/screenings?movie_id=${movieId}`);
        if (!response.ok) throw new Error('Failed to load screenings');

        screenings = await response.json();
        renderScreeningOptions();
        
        if (screenings.length > 0) {
            document.getElementById('screeningSelect').value = screenings[0].id;
        }
        
        document.getElementById('bookingModal').classList.add('show');
    } catch (error) {
        document.getElementById('bookingError').textContent = 'Failed to load screenings.';
        console.error('Load screenings error:', error);
    }
}

window.onclick = function(event) {
    const bookingModal = document.getElementById('bookingModal');
    if (event.target == bookingModal) {
        bookingModal.classList.remove('show');
    }
};
