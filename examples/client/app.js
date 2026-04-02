// ISM API Demo — Vanilla JS Client
// Zero dependencies. Demonstrates: guidance wizard, banner rendering,
// authority blocks, dissemination controls, validation, and more.

(function () {
  "use strict";

  // ---------------------------------------------------------------------------
  // State
  // ---------------------------------------------------------------------------

  var apiBase = "http://localhost:8080";
  var connected = false;

  // Cached reference data from the API.
  var refData = {
    classifications: [],
    cuiCategories: [],
    disseminationControls: [],
    distributionStatements: [],
    countryCodes: [],
    declassExceptions: [],
    nonICMarkings: [],
  };

  // Current ISM object being built.
  var ism = {};

  // Latest guidance response.
  var guidanceFields = [];

  // ---------------------------------------------------------------------------
  // DOM helpers
  // ---------------------------------------------------------------------------

  function $(id) {
    return document.getElementById(id);
  }

  function show(el) {
    if (typeof el === "string") el = $(el);
    if (el) el.classList.remove("hidden");
  }

  function hide(el) {
    if (typeof el === "string") el = $(el);
    if (el) el.classList.add("hidden");
  }

  // ---------------------------------------------------------------------------
  // API helpers
  // ---------------------------------------------------------------------------

  function apiGet(path) {
    return fetch(apiBase + path).then(function (res) {
      if (!res.ok) throw new Error("HTTP " + res.status);
      return res.json();
    });
  }

  function apiPost(path, body) {
    return fetch(apiBase + path, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(body),
    }).then(function (res) {
      if (!res.ok) throw new Error("HTTP " + res.status);
      return res.json();
    });
  }

  // ---------------------------------------------------------------------------
  // Connection & reference data loading
  // ---------------------------------------------------------------------------

  function connect() {
    apiBase = $("api-url").value.replace(/\/+$/, "");
    var status = $("api-status");
    status.textContent = "Connecting...";
    status.className = "";

    apiGet("/healthz")
      .then(function () {
        return loadRefData();
      })
      .then(function () {
        connected = true;
        status.textContent = "Connected";
        status.className = "status-connected";
        initWizard();
      })
      .catch(function (err) {
        connected = false;
        status.textContent = "Failed: " + err.message;
        status.className = "status-disconnected";
      });
  }

  function loadRefData() {
    return Promise.all([
      apiGet("/api/v1/ref/classifications"),
      apiGet("/api/v1/ref/cui-categories"),
      apiGet("/api/v1/ref/dissemination-controls"),
      apiGet("/api/v1/ref/distribution-statements"),
      apiGet("/api/v1/ref/country-codes"),
      apiGet("/api/v1/ref/declass-exceptions"),
      apiGet("/api/v1/ref/non-ic-markings"),
    ]).then(function (results) {
      refData.classifications = results[0].data || [];
      refData.cuiCategories = results[1].data || [];
      refData.disseminationControls = results[2].data || [];
      refData.distributionStatements = results[3].data || [];
      refData.countryCodes = results[4].data || [];
      refData.declassExceptions = results[5].data || [];
      refData.nonICMarkings = results[6].data || [];
    });
  }

  // ---------------------------------------------------------------------------
  // Wizard initialization
  // ---------------------------------------------------------------------------

  function initWizard() {
    populateClassifications();
    populateCountryCodes();
    populateDeclassExceptions();
    resetISM();
  }

  function populateClassifications() {
    var sel = $("field-classification");
    sel.innerHTML = '<option value="">— Select —</option>';
    refData.classifications.forEach(function (c) {
      var opt = document.createElement("option");
      opt.value = c.code;
      opt.textContent = c.code + " — " + c.label;
      sel.appendChild(opt);
    });
  }

  function populateCountryCodes() {
    var selectors = [
      "field-ownerProducer",
      "field-releasableTo",
      "field-displayOnlyTo",
      "field-fgiSourceOpen",
      "field-fgiSourceProtected",
    ];
    selectors.forEach(function (id) {
      var sel = $(id);
      sel.innerHTML = "";
      refData.countryCodes.forEach(function (c) {
        var opt = document.createElement("option");
        opt.value = c.code;
        opt.textContent = c.code + " — " + c.name;
        sel.appendChild(opt);
      });
    });
  }

  function populateDeclassExceptions() {
    var sel = $("field-declassException");
    sel.innerHTML = '<option value="">— None —</option>';
    refData.declassExceptions.forEach(function (c) {
      var opt = document.createElement("option");
      opt.value = c.code;
      opt.textContent = c.code + " — " + c.label;
      sel.appendChild(opt);
    });
  }

  // ---------------------------------------------------------------------------
  // ISM state management
  // ---------------------------------------------------------------------------

  function resetISM() {
    ism = {};
    guidanceFields = [];
    $("field-classification").value = "";
    clearMultiSelect("field-ownerProducer");
    clearMultiSelect("field-releasableTo");
    clearMultiSelect("field-displayOnlyTo");
    clearMultiSelect("field-fgiSourceOpen");
    clearMultiSelect("field-fgiSourceProtected");
    $("field-distributionStatement").value = "";
    $("field-classifiedBy").value = "";
    $("field-classificationReason").value = "";
    $("field-derivativelyClassifiedBy").value = "";
    $("field-derivedFrom").value = "";
    $("field-compilationReason").value = "";
    $("field-declassDate").value = "";
    $("field-declassEvent").value = "";
    $("field-declassException").value = "";
    hideAllSections();
    hide("validation-results");
    updatePreview();
  }

  function clearMultiSelect(id) {
    var sel = $(id);
    for (var i = 0; i < sel.options.length; i++) {
      sel.options[i].selected = false;
    }
  }

  function hideAllSections() {
    var sections = [
      "section-ownerProducer",
      "section-categoryMarkings",
      "section-disseminationControls",
      "section-releasableTo",
      "section-displayOnlyTo",
      "section-distributionStatement",
      "section-authority",
      "section-declass",
      "section-fgi",
      "section-nonICMarkings",
    ];
    sections.forEach(hide);
  }

  function buildISM() {
    ism = {};

    var cls = $("field-classification").value;
    if (cls) ism.classification = cls;

    // Owner/Producer
    var op = getMultiSelectValues("field-ownerProducer");
    if (op.length) ism.ownerProducer = op;

    // CUI Categories
    var cats = getCheckedValues("field-categoryMarkings");
    if (cats.length) ism.categoryMarkings = cats;

    // Dissemination Controls
    var dissem = getCheckedValues("field-disseminationControls");
    if (dissem.length) ism.disseminationControls = dissem;

    // Releasable To
    var rel = getMultiSelectValues("field-releasableTo");
    if (rel.length) ism.releasableTo = rel;

    // Display Only To
    var disp = getMultiSelectValues("field-displayOnlyTo");
    if (disp.length) ism.displayOnlyTo = disp;

    // Distribution Statement
    var dist = $("field-distributionStatement").value;
    if (dist) ism.distributionStatement = dist;

    // Authority
    var cb = $("field-classifiedBy").value.trim();
    if (cb) ism.classifiedBy = cb;

    var cr = $("field-classificationReason").value.trim();
    if (cr) ism.classificationReason = cr;

    var dcb = $("field-derivativelyClassifiedBy").value.trim();
    if (dcb) ism.derivativelyClassifiedBy = dcb;

    var df = $("field-derivedFrom").value.trim();
    if (df) ism.derivedFrom = df;

    var comp = $("field-compilationReason").value.trim();
    if (comp) ism.compilationReason = comp;

    // Declassification
    var dd = $("field-declassDate").value.trim();
    if (dd) ism.declassDate = dd;

    var de = $("field-declassEvent").value.trim();
    if (de) ism.declassEvent = de;

    var dx = $("field-declassException").value;
    if (dx) ism.declassException = dx;

    // FGI
    var fgiOpen = getMultiSelectValues("field-fgiSourceOpen");
    if (fgiOpen.length) ism.fgiSourceOpen = fgiOpen;

    var fgiProt = getMultiSelectValues("field-fgiSourceProtected");
    if (fgiProt.length) ism.fgiSourceProtected = fgiProt;

    // Non-IC Markings
    var nic = getCheckedValues("field-nonICMarkings");
    if (nic.length) ism.nonICMarkings = nic;

    return ism;
  }

  function getMultiSelectValues(id) {
    var sel = $(id);
    var vals = [];
    for (var i = 0; i < sel.options.length; i++) {
      if (sel.options[i].selected) vals.push(sel.options[i].value);
    }
    return vals;
  }

  function getCheckedValues(containerId) {
    var container = $(containerId);
    var vals = [];
    var boxes = container.querySelectorAll('input[type="checkbox"]:checked');
    for (var i = 0; i < boxes.length; i++) {
      vals.push(boxes[i].value);
    }
    return vals;
  }

  // ---------------------------------------------------------------------------
  // Guidance engine integration
  // ---------------------------------------------------------------------------

  function requestGuidance() {
    if (!connected) return;

    var currentISM = buildISM();
    apiPost("/api/v1/guidance", { ism: currentISM })
      .then(function (resp) {
        guidanceFields = (resp.data && resp.data.fields) || [];
        applyGuidance();
        updatePreview();
        $("guidance-json").textContent = JSON.stringify(
          guidanceFields,
          null,
          2
        );
      })
      .catch(function (err) {
        console.error("Guidance error:", err);
      });
  }

  function applyGuidance() {
    hideAllSections();

    // Build a lookup by field name.
    var byField = {};
    guidanceFields.forEach(function (fg) {
      byField[fg.field] = fg;
    });

    // Owner/Producer
    applyFieldSection("ownerProducer", byField);

    // CUI Categories
    var catGuide = byField["categoryMarkings"];
    if (catGuide && catGuide.status !== "not_applicable") {
      show("section-categoryMarkings");
      renderCheckboxes(
        "field-categoryMarkings",
        catGuide.allowedValues || [],
        getCheckedValues("field-categoryMarkings")
      );
      setHint("hint-categoryMarkings", catGuide);
    }

    // Dissemination Controls
    var dissemGuide = byField["disseminationControls"];
    if (dissemGuide && dissemGuide.status !== "not_applicable") {
      show("section-disseminationControls");
      renderCheckboxes(
        "field-disseminationControls",
        dissemGuide.allowedValues || [],
        getCheckedValues("field-disseminationControls")
      );
      setHint("hint-disseminationControls", dissemGuide);
    }

    // Releasable To
    applyFieldSection("releasableTo", byField);

    // Display Only To
    applyFieldSection("displayOnlyTo", byField);

    // Distribution Statement
    var distGuide = byField["distributionStatement"];
    if (distGuide && distGuide.status !== "not_applicable") {
      show("section-distributionStatement");
      populateDistributionStatements(distGuide.allowedValues || []);
      setHint("hint-distributionStatement", distGuide);
    }

    // Authority block (classifiedBy, etc.) — show for C/S
    var authGuide = byField["classifiedBy"];
    if (authGuide && authGuide.status !== "not_applicable") {
      show("section-authority");
      var authFields = [
        "classifiedBy",
        "classificationReason",
        "derivativelyClassifiedBy",
        "derivedFrom",
        "compilationReason",
      ];
      authFields.forEach(function (name) {
        applyFieldRequirement("field-" + name, byField[name]);
        applyFieldHint("fhint-" + name, byField[name]);
      });
    }

    // Declassification — show for C/S
    var declGuide = byField["declassDate"];
    if (declGuide && declGuide.status !== "not_applicable") {
      show("section-declass");
    }

    // FGI
    var fgiGuide = byField["fgiSourceOpen"];
    if (fgiGuide && fgiGuide.status !== "not_applicable") {
      show("section-fgi");
    }

    // Non-IC Markings
    var nicGuide = byField["nonICMarkings"];
    if (nicGuide && nicGuide.status !== "not_applicable") {
      show("section-nonICMarkings");
      renderCheckboxes(
        "field-nonICMarkings",
        nicGuide.allowedValues || [],
        getCheckedValues("field-nonICMarkings")
      );
    }
  }

  function applyFieldSection(fieldName, byField) {
    var guide = byField[fieldName];
    if (guide && guide.status !== "not_applicable") {
      show("section-" + fieldName);
      var hintEl = $("hint-" + fieldName);
      if (hintEl) setHint("hint-" + fieldName, guide);
    }
  }

  function applyFieldRequirement(inputId, guide) {
    if (!guide) return;
    var group = $(inputId).closest(".field-group");
    if (!group) return;
    if (guide.required || guide.status === "required") {
      group.classList.add("required");
    } else {
      group.classList.remove("required");
    }
  }

  function applyFieldHint(hintId, guide) {
    var el = $(hintId);
    if (!el) return;
    if (!guide) {
      el.textContent = "";
      return;
    }
    if (guide.required || guide.status === "required") {
      el.textContent = "Required";
    } else if (guide.requiredIf) {
      el.textContent = "Required if " + guide.requiredIf;
    } else {
      el.textContent = "";
    }
  }

  function setHint(hintId, guide) {
    var el = $(hintId);
    if (!el) return;
    var parts = [];
    if (guide.status) parts.push("Status: " + guide.status);
    if (guide.requiredIf) parts.push(guide.requiredIf);
    if (guide.reason) parts.push(guide.reason);
    el.textContent = parts.join(" — ");
  }

  function renderCheckboxes(containerId, allowedValues, checkedValues) {
    var container = $(containerId);
    container.innerHTML = "";
    allowedValues.forEach(function (av) {
      var label = document.createElement("label");
      var cb = document.createElement("input");
      cb.type = "checkbox";
      cb.value = av.code;
      cb.checked = checkedValues.indexOf(av.code) !== -1;
      cb.addEventListener("change", onFieldChange);
      label.appendChild(cb);
      var text = av.code;
      if (av.label && av.label !== av.code) text += " (" + av.label + ")";
      label.appendChild(document.createTextNode(" " + text));
      container.appendChild(label);
    });
  }

  function populateDistributionStatements(allowed) {
    var sel = $("field-distributionStatement");
    var current = sel.value;
    sel.innerHTML = '<option value="">— None —</option>';
    allowed.forEach(function (av) {
      var opt = document.createElement("option");
      opt.value = av.code;
      opt.textContent = av.code + " — " + av.label;
      sel.appendChild(opt);
    });
    sel.value = current;
  }

  // ---------------------------------------------------------------------------
  // Banner preview
  // ---------------------------------------------------------------------------

  function updatePreview() {
    var currentISM = buildISM();
    $("ism-json").textContent = JSON.stringify(currentISM, null, 2);

    if (!connected || !currentISM.classification) {
      $("banner-line").textContent = "UNCLASSIFIED";
      $("banner-line").className = "banner cls-U";
      $("portion-mark").textContent = "(U)";
      $("portion-mark").className = "portion cls-U";
      return;
    }

    apiPost("/api/v1/banner", { ism: currentISM })
      .then(function (resp) {
        var data = resp.data || {};
        $("banner-line").textContent = data.bannerLine || "—";
        $("portion-mark").textContent = data.portionMark || "—";

        var cls = currentISM.classification || "U";
        $("banner-line").className = "banner cls-" + cls;
        $("portion-mark").className = "portion cls-" + cls;
      })
      .catch(function (err) {
        $("banner-line").textContent = "Error: " + err.message;
        $("banner-line").className = "banner";
        $("portion-mark").textContent = "—";
        $("portion-mark").className = "portion";
      });
  }

  // ---------------------------------------------------------------------------
  // Validation
  // ---------------------------------------------------------------------------

  function validate() {
    if (!connected) return;

    var currentISM = buildISM();
    hide("validation-results");

    apiPost("/api/v1/validate", { ism: currentISM })
      .then(function (resp) {
        var data = resp.data || {};
        show("validation-results");

        var statusEl = $("validation-status");
        if (data.valid) {
          statusEl.textContent = "Valid";
          statusEl.className = "valid-true";
        } else {
          statusEl.textContent = "Invalid";
          statusEl.className = "valid-false";
        }

        var errList = $("validation-errors");
        errList.innerHTML = "";

        var errors = data.errors || [];
        if (errors.length === 0 && data.valid) {
          var li = document.createElement("li");
          li.textContent = "No errors or warnings.";
          li.style.color = "#2e7d32";
          errList.appendChild(li);
        }

        errors.forEach(function (e) {
          var li = document.createElement("li");
          li.className = "severity-" + e.severity;
          li.textContent =
            "[" +
            e.severity.toUpperCase() +
            "] " +
            e.field +
            ": " +
            e.message +
            " (" +
            e.code +
            ")";
          errList.appendChild(li);
        });
      })
      .catch(function (err) {
        show("validation-results");
        $("validation-status").textContent = "Request failed: " + err.message;
        $("validation-status").className = "valid-false";
        $("validation-errors").innerHTML = "";
      });
  }

  // ---------------------------------------------------------------------------
  // Event handlers
  // ---------------------------------------------------------------------------

  function onFieldChange() {
    requestGuidance();
  }

  function bindEvents() {
    // Connect button
    $("btn-connect").addEventListener("click", connect);
    $("api-url").addEventListener("keydown", function (e) {
      if (e.key === "Enter") connect();
    });

    // Classification dropdown — triggers full guidance refresh
    $("field-classification").addEventListener("change", onFieldChange);

    // Multi-selects
    $("field-ownerProducer").addEventListener("change", onFieldChange);
    $("field-releasableTo").addEventListener("change", onFieldChange);
    $("field-displayOnlyTo").addEventListener("change", onFieldChange);
    $("field-fgiSourceOpen").addEventListener("change", onFieldChange);
    $("field-fgiSourceProtected").addEventListener("change", onFieldChange);

    // Distribution statement
    $("field-distributionStatement").addEventListener("change", onFieldChange);

    // Authority text fields — debounced guidance on input
    var authorityFields = [
      "field-classifiedBy",
      "field-classificationReason",
      "field-derivativelyClassifiedBy",
      "field-derivedFrom",
      "field-compilationReason",
    ];
    authorityFields.forEach(function (id) {
      $(id).addEventListener("input", debounce(onFieldChange, 400));
    });

    // Declass fields
    $("field-declassDate").addEventListener(
      "input",
      debounce(onFieldChange, 400)
    );
    $("field-declassEvent").addEventListener(
      "input",
      debounce(onFieldChange, 400)
    );
    $("field-declassException").addEventListener("change", onFieldChange);

    // Validate and reset
    $("btn-validate").addEventListener("click", validate);
    $("btn-reset").addEventListener("click", function () {
      resetISM();
    });
  }

  function debounce(fn, ms) {
    var timer;
    return function () {
      clearTimeout(timer);
      timer = setTimeout(fn, ms);
    };
  }

  // ---------------------------------------------------------------------------
  // Boot
  // ---------------------------------------------------------------------------

  bindEvents();
  // Auto-connect on page load.
  connect();
})();
