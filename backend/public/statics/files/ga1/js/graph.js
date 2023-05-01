class Graph {
    width;
    height;
    root;
    treemap;
    treeData;
    duration = 700;

    svgElement= document.getElementById('svg');

    setup() {
        this.width = window.innerWidth;
        this.height = window.innerHeight;
        d3.select(this.svgElement)
            .attr("width", this.width + 1000)
            .attr("height", this.height - 20)
            .append("g")
            .attr("transform", function (_) {
                return "translate(30, 0)";
            })
        // declares a tree layout and assigns the size
        this.treemap = d3.tree().size([this.height, this.width]);
    }

    update(source) {
        let self = this;
        let i = 0;
        let svg = d3.select(this.svgElement).select('g');
        // Assigns the x and y position for the nodes

        // Compute the new tree layout.
        let nodes = this.treeData.descendants(),
            links = this.treeData.descendants().slice(1);

        // Normalize for fixed-depth.
        nodes.forEach(function (d) {
            d.y = d.depth * 180
        });

        // ****************** Nodes section ***************************

        // Update the nodes...
        let node = svg.selectAll('g.node')
            .data(nodes, function (d) {
                return d.id || (d.id = ++i);
            });

        // Enter any new modes at the parent's previous position.
        let nodeEnter = node.enter().append('g')
            .attr('class', 'node')
            .attr("transform", function (_) {
                return "translate(" + source.y0 + "," + source.x0 + ")";
            })
            .on('click', click);

        // Add Circle for the nodes
        nodeEnter.append('circle')
            .attr('class', 'node')
            .attr('r', 1e-6)
        // .style("fill", function (d) {
        //     return d._children ? "lightsteelblue" : "#fff";
        // });

        // Add labels for the nodes
        nodeEnter.append('text')
            .attr("dy", "2.2em")
            // .attr("dy", ".35em")
            // .attr("x", function (d) {
            //     return d.children || d._children ? -13 : 13;
            // })
            // .attr("text-anchor", function (d) {
            //     return d.children || d._children ? "end" : "start";
            // })
            .attr("stroke", "white")
            .text(function (d) {
                return d.data.name;
            });

        // UPDATE
        let nodeUpdate = nodeEnter.merge(node);

        // Transition to the proper position for the node
        nodeUpdate.transition()
            .duration(this.duration)
            .attr("transform", function (d) {
                return "translate(" + d.y + "," + d.x + ")";
            });

        // Update the node attributes and style
        nodeUpdate.select('circle.node')
            // .attr('r', 10)
            .attr('r', function (d) {
                return d._children ? 14 : 10;
            })
            .style("fill", function (d) {
                return d._children ? "lightsteelblue" : "#fff";
            })
            .attr('cursor', 'pointer');


        // Remove any exiting nodes
        let nodeExit = node.exit().transition()
            .duration(this.duration)
            .attr("transform", function (_) {
                return "translate(" + source.y + "," + source.x + ")";
            })
            .remove();

        // On exit reduce the node circles size to 0
        nodeExit.select('circle')
            .attr('r', 1e-6);

        // On exit reduce the opacity of text labels
        nodeExit.select('text')
            .style('fill-opacity', 1e-6);

        // ****************** links section ***************************

        // Update the links...
        let link = svg.selectAll('path.link')
            .data(links, function (d) {
                return d.id;
            });

        // Enter any new links at the parent's previous position.
        let linkEnter = link.enter().insert('path', "g")
            .attr("class", "link")
            .attr('d', function (_) {
                let o = {x: source.x0, y: source.y0}
                return diagonal(o, o)
            });

        // UPDATE
        let linkUpdate = linkEnter.merge(link);

        // Transition back to the parent element position
        linkUpdate.transition()
            .duration(this.duration)
            .attr('d', function (d) {
                return diagonal(d, d.parent)
            });

        // Remove any exiting links
        link.exit().transition()
            .duration(this.duration)
            .attr('d', function (_) {
                let o = {x: source.x, y: source.y}
                return diagonal(o, o)
            })
            .remove();

        // Store the old positions for transition.
        nodes.forEach(function (d) {
            d.x0 = d.x;
            d.y0 = d.y;
        });

        // Creates a curved (diagonal) path from parent to the child nodes
        function diagonal(s, d) {
            return `M ${s.y} ${s.x}
            C ${(s.y + d.y) / 2} ${s.x},
              ${(s.y + d.y) / 2} ${d.x},
              ${d.y} ${d.x}`
        }

        // Toggle children on click.
        function click(d) {
            if (d.children) {
                d._children = d.children;
                d.children = null;
            } else {
                d.children = d._children;
                d._children = null;
            }
            self.update(d);
        }
    }

    draw(data) {

        this.clear()

        // Assigns parent, children, height, depth
        this.root = d3.hierarchy(data, function (d) {
            return d.children;
        });
        this.root.x0 = this.height / 2;
        this.root.y0 = 0;

        this.treeData = d3.tree().size([this.height, this.width])(this.root);
        this.update(this.root);
    }

    clear() {
        // Select all nodes and links and remove them
        d3.select(this.svgElement)
          .selectAll('.node, .link')
          .remove();
    
        // Reset the tree data and call the update function
        this.treeData = null;
      }
}

