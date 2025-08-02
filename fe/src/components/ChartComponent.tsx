import React, { useEffect, useRef } from 'react';
import * as d3 from 'd3';

const ChartComponent: React.FC = () => {
  const svgRef = useRef<SVGSVGElement>(null);

  useEffect(() => {
    if (!svgRef.current) {return;}

    const svg = d3.select(svgRef.current);
    svg.selectAll('*').remove(); // Clear previous content

    const margin = { top: 20, right: 20, bottom: 30, left: 40 };
    const width = 400 - margin.left - margin.right;
    const height = 300 - margin.top - margin.bottom;

    // Sample data
    const data = [
      { name: 'A', value: 30 },
      { name: 'B', value: 80 },
      { name: 'C', value: 45 },
      { name: 'D', value: 60 },
      { name: 'E', value: 20 },
    ];

    // Create scales
    const xScale = d3.scaleBand()
      .domain(data.map(d => d.name))
      .range([0, width])
      .padding(0.1);

    const yScale = d3.scaleLinear()
      .domain([0, d3.max(data, d => d.value) || 0])
      .range([height, 0]);

    // Create main group
    const g = svg
      .attr('width', width + margin.left + margin.right)
      .attr('height', height + margin.top + margin.bottom)
      .append('g')
      .attr('transform', `translate(${margin.left},${margin.top})`);

    // Create bars
    g.selectAll('.bar')
      .data(data)
      .enter().append('rect')
      .attr('class', 'bar')
      .attr('x', d => xScale(d.name) || 0)
      .attr('width', xScale.bandwidth())
      .attr('y', d => yScale(d.value))
      .attr('height', d => height - yScale(d.value))
      .attr('fill', '#3B82F6')
      .on('mouseover', function(_event, _d) {
        d3.select(this).attr('fill', '#1D4ED8');
      })
      .on('mouseout', function(_event, _d) {
        d3.select(this).attr('fill', '#3B82F6');
      });

    // Add x-axis
    g.append('g')
      .attr('transform', `translate(0,${height})`)
      .call(d3.axisBottom(xScale));

    // Add y-axis
    g.append('g')
      .call(d3.axisLeft(yScale));

    // Add labels
    g.selectAll('.label')
      .data(data)
      .enter().append('text')
      .attr('class', 'label')
      .attr('x', d => (xScale(d.name) || 0) + xScale.bandwidth() / 2)
      .attr('y', d => yScale(d.value) - 5)
      .attr('text-anchor', 'middle')
      .attr('font-size', '12px')
      .attr('fill', '#374151')
      .text(d => d.value);

  }, []);

  return (
    <div className="w-full">
      <svg ref={svgRef} className="w-full h-auto"></svg>
    </div>
  );
};

export default ChartComponent;