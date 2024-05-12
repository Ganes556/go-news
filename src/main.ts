import './style/main.scss';
import * as bootstrap from 'bootstrap';
import Swal from 'sweetalert2';
import ckeditor from './ckeditor';
import htmx from './htmx';
// import * as dropzone from 'dropzone';
import Alpine from 'alpinejs';
import focus from '@alpinejs/focus';
import { formatDate } from './alpine_func';

declare global {
  interface Window {
    htmx: typeof htmx;
    ckeditor: typeof ckeditor;
    bootstrap: typeof bootstrap;
    Swal: typeof Swal;
    Alpine: typeof Alpine;
  }
}

window.htmx = htmx;
window.ckeditor = ckeditor;
window.bootstrap = bootstrap;
window.Swal = Swal;
window.Alpine = Alpine;

// formatting date
Alpine.data('date', () => ({
  formatDate,
}));

Alpine.plugin(focus);

Alpine.start();
