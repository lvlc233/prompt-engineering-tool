declare module 'react-graph-vis' {
  import { Component } from 'react';

  export interface Node {
    id: string | number;
    label?: string;
    title?: string;
    color?: {
      background?: string;
      border?: string;
      highlight?: {
        background?: string;
        border?: string;
      };
    };
    font?: {
      color?: string;
      size?: number;
    };
    size?: number;
    shape?: string;
    borderWidth?: number;
    shadow?: boolean;
  }

  export interface Edge {
    from: string | number;
    to: string | number;
    color?: {
      color?: string;
    };
    width?: number;
    arrows?: {
      to?: {
        enabled?: boolean;
        scaleFactor?: number;
      };
    };
  }

  export interface GraphData {
    nodes: Node[];
    edges: Edge[];
  }

  export interface GraphOptions {
    layout?: {
      hierarchical?: {
        direction?: string;
        sortMethod?: string;
        nodeSpacing?: number;
        levelSeparation?: number;
      };
    };
    edges?: {
      color?: string;
      smooth?: {
        type?: string;
        forceDirection?: string;
        roundness?: number;
      };
    };
    nodes?: {
      shape?: string;
      borderWidth?: number;
      shadow?: boolean;
    };
    physics?: {
      enabled?: boolean;
    };
    interaction?: {
      dragNodes?: boolean;
      dragView?: boolean;
      zoomView?: boolean;
    };
    height?: string;
    width?: string;
  }

  export interface GraphEvents {
    hoverNode?: (event: any) => void;
    blurNode?: (event: any) => void;
    click?: (event: any) => void;
    doubleClick?: (event: any) => void;
  }

  export interface GraphProps {
    graph: GraphData;
    options?: GraphOptions;
    events?: GraphEvents;
    getNetwork?: (network: any) => void;
    identifier?: string;
    style?: React.CSSProperties;
    getNodes?: (nodes: any) => void;
    getEdges?: (edges: any) => void;
  }

  export default class Graph extends Component<GraphProps> {}
}