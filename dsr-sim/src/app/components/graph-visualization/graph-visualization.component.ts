import {
  Component,
  ViewChild,
  ElementRef,
  OnDestroy,
  AfterViewInit,
} from '@angular/core';
import cytoscape from 'cytoscape';
import {
  GraphResponse,
  GraphService,
} from '../../services/graph.service.ts.service';
import { FormsModule } from '@angular/forms';
import { MatSliderDragEvent, MatSliderModule } from '@angular/material/slider';
import { MatButtonModule } from '@angular/material/button';
import { MatInputModule } from '@angular/material/input';

@Component({
  selector: 'app-graph-visualization',
  templateUrl: './graph-visualization.component.html',
  imports: [FormsModule, MatSliderModule, MatButtonModule, MatInputModule],
  styleUrls: ['./graph-visualization.component.css'],
})
export class GraphVisualizationComponent implements AfterViewInit, OnDestroy {
  @ViewChild('cy') cyContainer!: ElementRef;

  private cy: any;
  graphData: GraphResponse | null = null;

  nodesCount = 10;
  sourceId = 0;
  destId = 5;

  events: any[] = [];
  isRunning: boolean = false;

  constructor(private graphService: GraphService) {}

  ngAfterViewInit(): void {
    this.initCytoscape();
    this.loadGraph(this.nodesCount);
    this.resetVisualization();
  }

  private initCytoscape() {
    this.cy = cytoscape({
      container: this.cyContainer.nativeElement,
      style: [
        {
          selector: 'node',
          style: {
            'background-color': '#4a90e2',
            label: 'data(id)',
            width: 38,
            height: 38,
            'font-size': '20px',
            color: '#fff',
            'text-valign': 'center',
            'font-weight': 'bold',
          },
        },
        {
          selector: 'edge',
          style: {
            width: 2,
            'line-color': '#888',
            'curve-style': 'bezier',
          },
        },
      ],
      layout: { name: 'preset' },
      userZoomingEnabled: true,
      userPanningEnabled: true,
    });
  }

  loadGraph(nodesCount: number = 20) {
    this.resetVisualization();
    this.graphService.generateGraph(nodesCount).subscribe({
      next: (data) => {
        this.graphData = data;
        this.renderGraph(data);
      },
      error: (err) => console.error('Failed to load graph', err),
    });
  }

  nextStep() {
    this.graphService.nextStep().subscribe({
      next: (response) => {
        if (response.event) {
          this.handleSimulationEvent(response);
        }
      },
      error: (err) => console.error(err),
    });
  }

  private handleSimulationEvent(responseEvent: any) {
    this.events.unshift(responseEvent);
    const p = responseEvent.event.payload || {};
    const eventType = responseEvent.event.type.toLowerCase();

    const route = p.route_so_far || p.route || [];
    const currentNode = route.length > 0 ? route[route.length - 1] : p.from;
    const prevNode = route.length > 1 ? route[route.length - 2] : null;

    switch (eventType) {
      case 'rreq_processed':
      case 'rreq_propagated':
        this.highlightNode(currentNode, '#ffeb3b', 1000);
        if (prevNode !== null) {
          this.highlightEdge(prevNode, currentNode, '#60a5fa');
        }
        break;

      case 'rreq_dropped':
        const droppedNode =
          p.from !== undefined
            ? p.from
            : route.length > 0
              ? route[route.length - 1]
              : null;
        this.highlightNode(droppedNode, '#ef4444');
        this.highlightEdge(droppedNode, route[route.length - 2], '#ef4444');
        break;

      case 'rrep_generated':
        this.highlightPath(route, '#22c55e', 1500);
        break;

      case 'rrep_forwarded':
      case 'rrep_received':
        this.highlightNode(currentNode, '#4ade80', 1200);
        if (prevNode !== null) {
          this.highlightEdge(currentNode, prevNode, '#22c55e');
        }
        break;

      case 'route_discovered':
        this.highlightPath(route, '#eab308', 200000);
        break;
    }
  }

  private highlightPath(route: number[], color: string, time: number) {
    for (let i = 0; i < route.length - 1; i++) {
      this.highlightEdge(route[i], route[i + 1], color);
    }
    for (let i = 0; i < route.length; i++) {
      this.highlightNode(route[i], color, time);
    }
  }

  startRouteDiscovery() {
    if (!this.graphData) return;
    this.events = [];
    this.resetGraphColors();
    this.isRunning = true;

    this.graphService
      .startRouteDiscovery(this.sourceId, this.destId)
      .subscribe({
        next: () =>
          console.log(
            `Route discovery queued: ${this.sourceId} → ${this.destId}`,
          ),
        error: (err) => console.error(err),
      });
  }

  resetVisualization() {
    this.events = [];
    this.resetGraphColors();
  }

  private resetGraphColors() {
    if (!this.cy) return;
    this.cy.nodes().style('background-color', '#4a90e2');
    this.cy.edges().style({ 'line-color': '#888', width: 2 });
  }

  private renderGraph(graph: GraphResponse) {
    if (!this.cy) return;
    this.cy.elements().remove();
    const elements: any[] = [];

    graph.nodes.forEach((node) => {
      elements.push({
        data: { id: node.id.toString() },
        position: { x: node.x, y: node.y },
      });
    });

    graph.edges.forEach((edge) => {
      elements.push({
        data: {
          id: `e${edge.from}-${edge.to}`,
          source: edge.from.toString(),
          target: edge.to.toString(),
        },
      });
    });

    this.cy.add(elements);
    this.cy.layout({ name: 'preset', animate: true, duration: 500 }).run();
  }

  onSliderRelease(event: MatSliderDragEvent) {
    this.loadGraph(event.value);
  }

  highlightNode(
    nodeId: number,
    color: string = '#ffeb3b',
    duration: number = 800,
  ) {
    if (nodeId === undefined || nodeId === null) return;
    const node = this.cy.getElementById(nodeId.toString());
    if (node && node.length > 0) {
      const originalColor = node.style('background-color');
      node.style('background-color', color);
      setTimeout(() => node.style('background-color', originalColor), duration);
    }
  }

  highlightEdge(from: number, to: number, color: string = '#4ade80') {
    let edge = this.cy.getElementById(`e${from}-${to}`);
    if (!edge || edge.length === 0) {
      edge = this.cy.getElementById(`e${to}-${from}`);
    }

    if (edge && edge.length > 0) {
      edge.style({ 'line-color': color, width: 4 });
    }
  }

  ngOnDestroy() {
    if (this.cy) this.cy.destroy();
  }
}
