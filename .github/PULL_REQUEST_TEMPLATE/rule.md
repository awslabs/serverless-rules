---
name: New Rule
about: Pull Request template for new rules
title: "rule: "
labels: new-rule
assignees: ""
---

## Key information

* Rule issue: _(add link to issue)_
* `cfn-lint` rule ID: 
* `tflint` rule name: 

## Summary

> One paragraph explaining the proposed rule.

## Checklist

* [ ] __cfn-lint__: add rule
* [ ] __cfn-lint__: add import in `rules/__init__.py`
* [ ] __cfn-lint__: add test templates
* [ ] __tflint__: add rule
* [ ] __tflint__: add tests
* [ ] __tflint__: add in `rules/provider.go`
* [ ] __docs__: add new section in documentation
* [ ] __docs__: add reference in `rules/index.md`