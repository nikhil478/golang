package images

const ArchitecturePrompt = `
You are an expert system architecture reverse engineer.

Analyze the image and extract a dependency graph.

Return JSON ONLY.

OUTPUT FORMAT:

{
  "nodes": [
    {
      "id": "service-name",
      "type": "service|db|queue|api|unknown"
    }
  ],
  "edges": [
    {
      "from": "node-id",
      "to": "node-id",
      "label": "optional description"
    }
  ]
}

RULES:
- Extract ONLY what is visible in the diagram
- DO NOT hallucinate nodes or connections
- If uncertain, omit edges
- Preserve exact naming from the image
- Focus on arrows, flows, and connections
`