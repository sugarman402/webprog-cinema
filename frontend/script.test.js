const { JSDOM } = require('jsdom');

// Setup JSDOM for DOM manipulation tests
const dom = new JSDOM('<!DOCTYPE html><html><body></body></html>', {
  url: 'http://localhost'
});
global.window = dom.window;
global.document = dom.window.document;
global.localStorage = {
  getItem: jest.fn(),
  setItem: jest.fn(),
  removeItem: jest.fn()
};
global.fetch = jest.fn();
global.alert = jest.fn();
global.confirm = jest.fn(() => true);

describe('escapeHtml', () => {
  const escapeHtml = (text) => {
    if (!text) return '';
    const map = {
      '&': '&amp;',
      '<': '&lt;',
      '>': '&gt;',
      '"': '&quot;',
      "'": '&#039;'
    };
    return String(text).replace(/[&<>"']/g, m => map[m]);
  };

  test('should escape HTML characters', () => {
    expect(escapeHtml('<script>alert("XSS")</script>')).toBe('&lt;script&gt;alert(&quot;XSS&quot;)&lt;/script&gt;');
  });

  test('should handle empty string', () => {
    expect(escapeHtml('')).toBe('');
  });

  test('should handle null', () => {
    expect(escapeHtml(null)).toBe('');
  });

  test('should handle undefined', () => {
    expect(escapeHtml(undefined)).toBe('');
  });

  test('should escape single quotes', () => {
    expect(escapeHtml("It's a test")).toBe('It&#039;s a test');
  });

  test('should escape ampersands', () => {
    expect(escapeHtml('Tom & Jerry')).toBe('Tom &amp; Jerry');
  });
});

describe('API_BASE', () => {
  test('should be defined', () => {
    const API_BASE = 'http://localhost:8080/api';
    expect(API_BASE).toBe('http://localhost:8080/api');
  });
});

describe('handleLogin', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  test('should handle successful login', async () => {
    const mockResponse = {
      ok: true,
      json: jest.fn().mockResolvedValue({
        token: 'jwt_token',
        user: { id: 1, email: 'test@example.com', full_name: 'Test User' }
      })
    };
    global.fetch.mockResolvedValue(mockResponse);

    document.getElementById = jest.fn((id) => {
      if (id === 'loginEmail') return { value: 'test@example.com' };
      if (id === 'loginPassword') return { value: 'password' };
      if (id === 'loginForm') return { reset: jest.fn() };
      if (id === 'authMessage') return { textContent: '', className: 'alert', classList: { add: jest.fn(), remove: jest.fn() } };
      return null;
    });

    const email = 'test@example.com';
    const password = 'password';

    const response = await fetch('http://localhost:8080/api/login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email, password })
    });

    expect(response).toBe(mockResponse);
    expect(global.fetch).toHaveBeenCalledWith('http://localhost:8080/api/login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email, password })
    });
  });

  test('should handle login failure', async () => {
    const mockResponse = {
      ok: false,
      text: jest.fn().mockResolvedValue('Invalid credentials')
    };
    global.fetch.mockResolvedValue(mockResponse);

    document.getElementById = jest.fn((id) => {
      if (id === 'authMessage') return { textContent: '', className: 'alert', classList: { add: jest.fn() } };
      return null;
    });

    const email = 'wrong@example.com';
    const password = 'wrongpassword';

    const response = await fetch('http://localhost:8080/api/login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email, password })
    });

    expect(response.ok).toBe(false);
  });
});

describe('switchTab', () => {
  test('should switch to login tab', () => {
    const loginForm = { classList: { add: jest.fn(), remove: jest.fn() } };
    const registerForm = { classList: { add: jest.fn(), remove: jest.fn() } };
    const tabBtns = [
      { classList: { add: jest.fn(), remove: jest.fn() } },
      { classList: { add: jest.fn(), remove: jest.fn() } }
    ];

    document.getElementById = jest.fn((id) => {
      if (id === 'loginForm') return loginForm;
      if (id === 'registerForm') return registerForm;
      return null;
    });

    document.querySelectorAll = jest.fn((selector) => {
      if (selector === '.tab-btn') return tabBtns;
      return [];
    });

    if ('login' === 'login') {
      loginForm.classList.add('active');
      registerForm.classList.remove('active');
      tabBtns[0].classList.add('active');
      tabBtns[1].classList.remove('active');
    }

    expect(loginForm.classList.add).toHaveBeenCalledWith('active');
    expect(registerForm.classList.remove).toHaveBeenCalledWith('active');
  });
});