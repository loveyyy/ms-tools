/*
 * @Author: Morningsun 1049935747@qq.com
 * @Date: 2025-09-07 17:36:40
 * @LastEditors: Morningsun 1049935747@qq.com
 * @LastEditTime: 2025-09-07 17:37:14
 * @Description: 
 */

import { marked } from 'marked'

import hljs from 'highlight.js'
import DOMPurify from 'dompurify'
import 'highlight.js/styles/github.css'

import katex from "katex";
import "katex/dist/katex.min.css";


export function renderMarkdown (content) {
    content = content.replace(/\$\$(.+?)\$\$/gs, (_, equation) => {
      return (
        "<katex-formula-ml>" +
        katex.renderToString(equation, {
          throwOnError: true,
          displayMode: true,
        }) +
        "</katex-formula-ml>"
      );
    }).replace(/\$(.+?)\$/g, (_, equation) => {
      return (
        "<katex-formula>" +
        katex.renderToString(equation, {
          throwOnError: true,
          displayMode: false,
        }) +
        "</katex-formula>"
      );
    });
  
    marked.setOptions({
      highlight: (code, lang) => {
        const language = hljs.getLanguage(lang) ? lang : 'plaintext';
        return hljs.highlight(code, { language }).value;
      },
      breaks: true
    });
  
    return DOMPurify.sanitize(marked(content));
  }