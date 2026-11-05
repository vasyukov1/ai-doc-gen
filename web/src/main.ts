import "./styles/main.css";
import { initAssembly } from "./assembly";
import { initForms } from "./forms";
import { initJobPolling } from "./job";
import { initPdfViewer } from "./pdf";
import { initTextareas } from "./textareas";
import { initTheme } from "./theme";

initTheme();
initForms();
initTextareas();
initAssembly();
initJobPolling();
initPdfViewer();

import { initSetup } from "./setup";
import { initUsage } from "./usage";
initSetup();
initUsage();

import { initSectionsReading } from "./sections";
initSectionsReading();
