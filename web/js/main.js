// Scroll reveal functionality
document.addEventListener('DOMContentLoaded', () => {
    const observerOptions = {
        root: null,
        rootMargin: '0px',
        threshold: 0.1
    };

    const observer = new IntersectionObserver((entries) => {
        entries.forEach(entry => {
            if (entry.isIntersecting) {
                entry.target.classList.add('is-visible');
            }
        });
    }, observerOptions);

    // Observe all elements with reveal-on-scroll class
    document.querySelectorAll('.reveal-on-scroll').forEach((element) => {
        observer.observe(element);
    });

    // Smooth scroll for navigation links
    document.querySelectorAll('a[href^="#"]').forEach(anchor => {
        anchor.addEventListener('click', function (e) {
            e.preventDefault();
            const target = document.querySelector(this.getAttribute('href'));
            if (target) {
                target.scrollIntoView({
                    behavior: 'smooth'
                });
            }
        });
    });

    // Add hover effects to interactive elements
    document.querySelectorAll('.bg-white.p-6').forEach(card => {
        card.classList.add('hover-scale');
    });
});

// Sample email data
const emails = [
    {
        id: 1,
        sender: "John Doe",
        subject: "Project Update",
        preview: "Here are the latest updates on the project timeline and milestones...",
        folder: "inbox",
        unread: true,
        date: "2025-11-10"
    },
    {
        id: 2,
        sender: "Alice Smith",
        subject: "Meeting Schedule",
        preview: "Can we schedule a meeting for next week to discuss the new features?",
        folder: "inbox",
        unread: true,
        date: "2025-11-09"
    },
    {
        id: 3,
        sender: "Marketing Team",
        subject: "Q4 Newsletter",
        preview: "Please review the draft of our Q4 newsletter before we send it out...",
        folder: "sent",
        unread: false,
        date: "2025-11-08"
    },
    {
        id: 4,
        sender: "Support System",
        subject: "Your Account Security",
        preview: "We noticed a login attempt from a new device...",
        folder: "spam",
        unread: true,
        date: "2025-11-07"
    },
    {
        id: 5,
        sender: "David Wilson",
        subject: "Holiday Party Planning",
        preview: "Let's start planning for the annual holiday party...",
        folder: "inbox",
        unread: false,
        date: "2025-11-06"
    },
    {
        id: 6,
        sender: "Sarah Johnson",
        subject: "Code Review Request",
        preview: "Could you please review my latest pull request when you have a moment?",
        folder: "sent",
        unread: false,
        date: "2025-11-05"
    }
];

// Current state
let currentFolder = 'inbox';
let currentEmails = [...emails];

// DOM Elements
document.addEventListener('DOMContentLoaded', () => {
    const loginForm = document.getElementById('loginForm');
    const loginScreen = document.getElementById('loginScreen');
    const emailInterface = document.getElementById('emailInterface');
    const logoutBtn = document.getElementById('logoutBtn');
    const searchInput = document.getElementById('searchEmails');
    const unreadOnlyCheckbox = document.getElementById('unreadOnly');
    const emailList = document.getElementById('emailList');
    const folderButtons = document.querySelectorAll('.mail-folder');

    // Login form handler
    loginForm.addEventListener('submit', (e) => {
        e.preventDefault();
        loginScreen.classList.add('hidden');
        emailInterface.classList.remove('hidden');
        emailInterface.classList.add('animate-fade-in');
        renderEmails();
    });

    // Logout handler
    logoutBtn.addEventListener('click', () => {
        emailInterface.classList.add('hidden');
        loginScreen.classList.remove('hidden');
        loginScreen.classList.add('animate-fade-in');
    });

    // Folder selection
    folderButtons.forEach(button => {
        button.addEventListener('click', () => {
            // Update active state
            folderButtons.forEach(btn => btn.classList.remove('active'));
            button.classList.add('active');

            // Update current folder and render
            currentFolder = button.dataset.folder;
            filterAndRenderEmails();
        });
    });

    // Search handler
    searchInput.addEventListener('input', filterAndRenderEmails);

    // Unread filter handler
    unreadOnlyCheckbox.addEventListener('change', filterAndRenderEmails);

    // Initial render
    document.querySelector('[data-folder="inbox"]').classList.add('active');
});

// Filter and render emails
function filterAndRenderEmails() {
    const searchTerm = document.getElementById('searchEmails').value.toLowerCase();
    const unreadOnly = document.getElementById('unreadOnly').checked;

    currentEmails = emails.filter(email => {
        const matchesFolder = email.folder === currentFolder;
        const matchesSearch = email.subject.toLowerCase().includes(searchTerm) ||
                            email.sender.toLowerCase().includes(searchTerm) ||
                            email.preview.toLowerCase().includes(searchTerm);
        const matchesUnread = !unreadOnly || email.unread;

        return matchesFolder && matchesSearch && matchesUnread;
    });

    renderEmails();
}

// Render emails
function renderEmails() {
    const emailList = document.getElementById('emailList');
    emailList.innerHTML = '';

    currentEmails.forEach(email => {
        const emailElement = document.createElement('div');
        emailElement.className = `email-card bg-white p-4 rounded-lg ${email.unread ? 'unread' : ''} animate-fade-in`;
        
        emailElement.innerHTML = `
            <div class="flex justify-between items-start">
                <div>
                    <h3 class="text-lg font-semibold">${email.sender}</h3>
                    <h4 class="text-md text-gray-700">${email.subject}</h4>
                    <p class="text-gray-600 mt-2">${email.preview}</p>
                </div>
                <div class="text-sm text-gray-400">
                    ${email.date}
                    ${email.unread ? '<span class="ml-2 bg-blue-100 text-blue-800 text-xs px-2 py-1 rounded">New</span>' : ''}
                </div>
            </div>
        `;

        // Mark as read on click
        emailElement.addEventListener('click', () => {
            if (email.unread) {
                email.unread = false;
                emailElement.classList.remove('unread');
                filterAndRenderEmails();
            }
        });

        emailList.appendChild(emailElement);
    });
}