import './style/main.scss';
import * as bootstrap from 'bootstrap';
import Swal from 'sweetalert2';
import ckeditor from './ckeditor';
import htmx from './htmx';
// import * as dropzone from 'dropzone';

declare global {
  interface Window {
    htmx: typeof htmx;
    ckeditor: typeof ckeditor;
    bootstrap: typeof bootstrap;
    Swal: typeof Swal;
    // Dropzone: typeof dropzone;
  }
}

window.htmx = htmx;
window.ckeditor = ckeditor;
window.bootstrap = bootstrap;
window.Swal = Swal;
// window.Dropzone = dropzone;
