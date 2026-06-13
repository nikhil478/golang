package images


// Step 5: Image Analyzer

// Use a multimodal model.

// Recommended:

// Gemini
// GPT-4o
// Claude

// Pipeline:

// PNG
//   |
// Vision Model
//   |
// JSON Facts

// Prompt:

// Extract:

// - services
// - databases
// - queues
// - APIs
// - arrows
// - dependencies

// Return JSON only.

// Output:

// {
//   "services":[
//     "Payment Service",
//     "Ledger Service"
//   ],
//   "dependencies":[
//     {
//       "from":"Payment Service",
//       "to":"Kafka"
//     }
//   ]
// }

// Convert everything to Facts.