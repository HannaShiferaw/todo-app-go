// API Base URL
const API_BASE = '/api'; // same origin, so it works on any host and port

// DOM Elements
const todoTitle = document.getElementById('todoTitle');
const todoDescription = document.getElementById('todoDescription');
const todoPriority = document.getElementById('todoPriority');
const todoCategory = document.getElementById('todoCategory');
const todoDueDate = document.getElementById('todoDueDate');
const addBtn = document.getElementById('addBtn');
const todosList = document.getElementById('todosList');
const emptyState = document.getElementById('emptyState');
const filterBtns = document.querySelectorAll('.filter-btn');
const totalTasksEl = document.getElementById('totalTasks');
const completedTasksEl = document.getElementById('completedTasks');
const remainingTasksEl = document.getElementById('remainingTasks');


let todos = [];
let currentFilter = 'all';


document.addEventListener('DOMContentLoaded', () => {
    loadTodos();
    setupEventListeners();
});


function setupEventListeners() {
    addBtn.addEventListener('click', addTodo);
    todoTitle.addEventListener('keypress', (e) => {
        if (e.key === 'Enter') addTodo();
    });

    filterBtns.forEach(btn => {
        btn.addEventListener('click', () => {
            filterBtns.forEach(b => b.classList.remove('active'));
            btn.classList.add('active');
            currentFilter = btn.dataset.filter;
            renderTodos();
        });
    });
}


async function loadTodos() {
    try {
        addBtn.innerHTML = '<i class="fas fa-spinner spinner"></i> Loading...';
        addBtn.disabled = true;

        const response = await fetch(`${API_BASE}/todos`);
        if (!response.ok) throw new Error('Failed to load todos');

        todos = await response.json() || [];
        renderTodos();
        updateStats();

        addBtn.innerHTML = '<i class="fas fa-plus"></i> Add Task';
        addBtn.disabled = false;
    } catch (error) {
        console.error('Error loading todos:', error);
        showToast('Failed to load todos', 'error');
        addBtn.innerHTML = '<i class="fas fa-plus"></i> Add Task';
        addBtn.disabled = false;
    }
}


async function addTodo() {
    const title = todoTitle.value.trim();
    const description = todoDescription.value.trim();
    const priority = todoPriority.value;
    const category = todoCategory.value;
    const dueDate = todoDueDate.value;

    if (!title) {
        showToast('Please enter a task title', 'error');
        todoTitle.focus();
        return;
    }

    const newTodo = {
        title,
        description,
        priority,
        category,
        completed: false,
        due_date: dueDate ? new Date(dueDate).toISOString() : null
    };

    try {
        addBtn.innerHTML = '<i class="fas fa-spinner spinner"></i> Adding...';
        addBtn.disabled = true;

        const response = await fetch(`${API_BASE}/todos`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(newTodo)
        });

        if (!response.ok) throw new Error('Failed to add todo');

        const createdTodo = await response.json();
        todos.push(createdTodo);

        
        todoTitle.value = '';
        todoDescription.value = '';
        todoPriority.value = 'medium';
        todoCategory.value = 'personal';
        todoDueDate.value = '';

        renderTodos();
        updateStats();
        showToast('Task added successfully!', 'success');

        addBtn.innerHTML = '<i class="fas fa-plus"></i> Add Task';
        addBtn.disabled = false;
    } catch (error) {
        console.error('Error adding todo:', error);
        showToast('Failed to add task', 'error');
        addBtn.innerHTML = '<i class="fas fa-plus"></i> Add Task';
        addBtn.disabled = false;
    }
}


async function updateTodo(id, updates) {
    try {
        const response = await fetch(`${API_BASE}/todos/${id}`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(updates)
        });

        if (!response.ok) throw new Error('Failed to update todo');

        const updatedTodo = await response.json();
        const index = todos.findIndex(t => t.id === id);
        if (index !== -1) {
            todos[index] = updatedTodo;
        }

        renderTodos();
        updateStats();
        showToast('Task updated successfully!', 'success');
    } catch (error) {
        console.error('Error updating todo:', error);
        showToast('Failed to update task', 'error');
        loadTodos(); 
    }
}


async function deleteTodo(id) {
    if (!confirm('Are you sure you want to delete this task?')) return;

    try {
        const response = await fetch(`${API_BASE}/todos/${id}`, {
            method: 'DELETE'
        });

        if (!response.ok) throw new Error('Failed to delete todo');

        todos = todos.filter(t => t.id !== id);
        renderTodos();
        updateStats();
        showToast('Task deleted successfully!', 'success');
    } catch (error) {
        console.error('Error deleting todo:', error);
        showToast('Failed to delete task', 'error');
        loadTodos(); 
    }
}


async function toggleTodoCompletion(id) {
    const todo = todos.find(t => t.id === id);
    if (!todo) return;

    await updateTodo(id, {
        ...todo,
        completed: !todo.completed
    });
}


function getFilteredTodos() {
    switch (currentFilter) {
        case 'active':
            return todos.filter(t => !t.completed);
        case 'completed':
            return todos.filter(t => t.completed);
        case 'high':
            return todos.filter(t => t.priority === 'high');
        case 'work':
            return todos.filter(t => t.category === 'work');
        default:
            return todos;
    }
}


function renderTodos() {
    const filteredTodos = getFilteredTodos();

    if (filteredTodos.length === 0) {
        todosList.innerHTML = '';
        emptyState.style.display = 'block';
        return;
    }

    emptyState.style.display = 'none';
    todosList.innerHTML = filteredTodos.map(todo => createTodoElement(todo)).join('');

   
    document.querySelectorAll('.todo-checkbox').forEach(checkbox => {
        checkbox.addEventListener('change', (e) => {
            toggleTodoCompletion(parseInt(e.target.dataset.id));
        });
    });

    document.querySelectorAll('.btn-delete').forEach(btn => {
        btn.addEventListener('click', (e) => {
            deleteTodo(parseInt(e.target.dataset.id));
        });
    });
}


function createTodoElement(todo) {
    const dueDate = todo.due_date ? new Date(todo.due_date) : null;
    const dueDateStr = dueDate ? formatDate(dueDate) : '';
    const isOverdue = dueDate && dueDate < new Date() && !todo.completed;

    return `
        <div class="todo-item ${todo.completed ? 'completed' : ''} ${todo.priority === 'high' ? 'high-priority' : ''}">
            <input 
                type="checkbox" 
                class="todo-checkbox" 
                ${todo.completed ? 'checked' : ''}
                data-id="${todo.id}"
            >
            <div class="todo-content">
                <div class="todo-text">${escapeHtml(todo.title)}</div>
                ${todo.description ? `<div class="todo-description">${escapeHtml(todo.description)}</div>` : ''}
                <div class="todo-meta">
                    <span class="priority-badge ${todo.priority}">${capitalizeFirst(todo.priority)} Priority</span>
                    <span class="category-badge">${capitalizeFirst(todo.category)}</span>
                    ${dueDateStr ? `<span class="due-date-badge ${isOverdue ? 'overdue' : ''}">${isOverdue ? '⚠️ ' : '📅 '}${dueDateStr}</span>` : ''}
                </div>
            </div>
            <div class="todo-actions">
                <button class="btn-small btn-delete" data-id="${todo.id}">
                    <i class="fas fa-trash"></i> Delete
                </button>
            </div>
        </div>
    `;
}


function updateStats() {
    const total = todos.length;
    const completed = todos.filter(t => t.completed).length;
    const remaining = total - completed;

    totalTasksEl.textContent = total;
    completedTasksEl.textContent = completed;
    remainingTasksEl.textContent = remaining;
}

// Utility functions
function formatDate(date) {
    const today = new Date();
    const tomorrow = new Date(today);
    tomorrow.setDate(tomorrow.getDate() + 1);

    if (date.toDateString() === today.toDateString()) {
        return 'Today';
    } else if (date.toDateString() === tomorrow.toDateString()) {
        return 'Tomorrow';
    } else {
        return date.toLocaleDateString('en-US', { month: 'short', day: 'numeric' });
    }
}

function capitalizeFirst(str) {
    return str.charAt(0).toUpperCase() + str.slice(1);
}

function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

function showToast(message, type = 'info') {
    const toast = document.createElement('div');
    toast.className = 'toast';
    toast.style.background = type === 'success' ? '#10b981' : type === 'error' ? '#ef4444' : '#6366f1';
    toast.textContent = message;
    document.body.appendChild(toast);

    setTimeout(() => {
        toast.style.opacity = '0';
        toast.style.transform = 'translateY(20px)';
        setTimeout(() => toast.remove(), 300);
    }, 3000);
}


window.addEventListener('error', (e) => {
    console.error('Unexpected error:', e);
});
